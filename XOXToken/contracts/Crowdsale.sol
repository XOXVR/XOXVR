pragma solidity ^0.4.23;

import './SafeMath.sol';
import './XOXToken.sol';
import './Ownable.sol';

contract Crowdsale is Ownable {
  using SafeMath for uint;

  uint constant public SOFT_CAP = 1E24;
  uint constant private DECIMALS = 1E18;
  
  uint public paginationNum = 100;

  uint public countRefund;
  uint public rate = 136900000000000;
  uint public start;
  uint public salesTokens;
  
  bool public requireOnce = true;

  address[] public investors;
  address public multisig;
  address public proxyMultisig;

  mapping(address => uint) public investments;
  
  XOXToken public token;

  struct Stage {
    uint id;
    uint totalSupply;
    uint remainderTokens;
    uint minInvestment;
    uint sumInvestment;
    uint end;
    uint bonusLessTen;
    uint bonusMoreTen;
    uint start;
  }

    Stage public preICO;
    Stage public ICO;
    Stage public postICO;


  event TransferWei(address indexed addr, uint amount);

  constructor(uint _startPreICO, uint _startICO, uint _startPostICO, address _reserved) public {
    require(
      _startPreICO > now &&
      _startICO > _startPreICO.add(28 days) &&
      _startPostICO > _startICO.add(28 days)
    );
    require(requireOnce);

    requireOnce = false;
    multisig = owner;
    token = new XOXToken(_reserved);
    start = _startPreICO;

    preICO = Stage({
      id: 0,
      totalSupply: 6E24,
      remainderTokens: 6E24,
      minInvestment: 1 ether,
      sumInvestment: 0,
      end: start.add(28 days),
      bonusLessTen: 140,
      bonusMoreTen: 160,
      start: _startPreICO
    });

    ICO = Stage({
      id: 1,
      totalSupply: 54E24,
      remainderTokens: 54E24,
      minInvestment: 0.5 ether,
      sumInvestment: 0,
      end: _startICO.add(28 days),
      bonusLessTen: 0,
      bonusMoreTen: 140,
      start: _startICO
    });

    postICO = Stage({
      id: 2,
      totalSupply: 840E24,
      remainderTokens: 840E24,
      minInvestment: 0,
      sumInvestment: 0,
      end: 0,
      bonusLessTen: 0,
      bonusMoreTen: 0,
      start: _startPostICO
    });
  }

  modifier onlyOwnerAndProxyOwner() {
    require(msg.sender == owner || msg.sender == proxyMultisig);
    _;
  }
  
  
  function() public payable {
    saleTokens();
  }
  

  function saleTokens() public payable {
    Stage storage stageICO = getStage();
    require(stageICO.id != 2);
    require(msg.value >= stageICO.minInvestment);
   
    uint tokensValue = calculateTokens(msg.value);
    
    if (tokensValue <= stageICO.remainderTokens) {
      uint tokens = bonusSystem(msg.value, stageICO, tokensValue);
      stageICO.remainderTokens = stageICO.remainderTokens.sub(tokens);
      addStatisticsAndCalculateVal(tokens, msg.value, msg.sender);
    } else {
        calculateTokensLessAvailable(stageICO, msg.value, msg.sender);
    }
    withDrawal();
  }


  function bonusSystem(uint _eth, Stage storage _stageNow, uint _tokens) internal view returns(uint) {
    uint eth = _eth.div(1E18);
    uint tokens = _tokens;
    Stage storage stage = _stageNow;
    if (stage.id == 0) {
        if (eth < 10) {
          tokens = tokens.mul(stage.bonusLessTen).div(100);
        } else if (eth >= 10) {
          tokens = tokens.mul(stage.bonusMoreTen).div(100);
        }
    } else {
       if (eth < 10) {
          tokens = tokens;
       } else if (eth >= 10) {
          tokens = tokens.mul(stage.bonusMoreTen).div(100);
       }
    }
    return tokens;
  }


  function getValueForBackend(uint _btc, address _sender) public onlyOwnerAndProxyOwner returns(uint) {
    require(_btc > 0);
    require(_sender != address(0));
    require(now > ICO.end);
    
    uint tokensValue = calculateTokens(_btc);
    uint surrender;
    if (tokensValue <= postICO.remainderTokens) {
      postICO.remainderTokens = postICO.remainderTokens.sub(tokensValue);
      surrender = 0;
    } else {
        tokensValue = postICO.remainderTokens;
        postICO.remainderTokens = 0;
        uint availableWei = tokensValue.mul(rate).div(DECIMALS);
        uint remainderWei = _btc.sub(availableWei);
        surrender = remainderWei;
    }
    token.transfer(_sender, tokensValue);
    return surrender;
  }


  function calculateTokensLessAvailable(Stage storage _stageNow, uint _ethValue, address _sender) internal {
      Stage storage stage = _stageNow;
      uint tokensValue = stage.remainderTokens;
      stage.remainderTokens = 0;
      uint tokens = bonusSystem(msg.value, stage, tokensValue);
      uint availableWei = tokens.mul(rate).div(DECIMALS);
      uint remainderWei = _ethValue.sub(availableWei);
    
      addStatisticsAndCalculateVal(tokens, availableWei, _sender);

      if (remainderWei > 0) {
        emit TransferWei(_sender, remainderWei);
        _sender.transfer(remainderWei);
      }
  }


  function addStatisticsAndCalculateVal(uint _tokensVal, uint _ethVal, address _sender) internal {
    if (investments[_sender] == 0) investors.push(_sender);
      investments[_sender] = investments[_sender].add(_ethVal);
      salesTokens = salesTokens.add(_tokensVal);
      token.transfer(_sender, _tokensVal);
  }


  function getStage() internal view returns(Stage storage) {
    if (now < preICO.end && now >= preICO.start) {
        return preICO;
    } else if (ICO.end > now && now >= ICO.start) {
        return ICO;
    } else if (now >= postICO.start) {
        return postICO;
    } else {
        revert();
    }
  }


  function withDrawal() internal {
    if(salesTokens >= SOFT_CAP) {
      multisig.transfer(address(this).balance);
    }
  }

  function burnTokens() public onlyOwnerAndProxyOwner {
    require(ICO.end < now);
    uint remainderTokens = preICO.remainderTokens.add(ICO.remainderTokens);
    token.burn(this, remainderTokens);
  }


  function calculateTokens(uint _ether) public view returns(uint) {
    return _ether.mul(DECIMALS).div(rate);
  }
    

  function calculateTotal(uint _tokens) public view returns (uint) {
    return _tokens.mul(rate);
  }
  
  
  function returnNumberInvestors() public view returns(uint) {
    return investors.length;
  }
    

  function refund() public onlyOwnerAndProxyOwner {
    require(ICO.end < now && salesTokens < SOFT_CAP);
    uint limit = countRefund.add(paginationNum);

    for(uint i = countRefund; i < investors.length; i++) {
      uint value = investments[investors[i]];
      investments[investors[i]] = 0;
      investors[i].transfer(value);
      if (i == limit) break;
    }
    countRefund = limit;
  }
  

  function setRate(uint _rate) public onlyOwner {
    require(_rate > 0);
    rate = _rate.div(10);
  }


  function setMultisig(address _to) public onlyOwner returns(bool) {
    require(_to != address(0));
    multisig = _to;
    return true;
  }

  function setDate(uint _startPreICO, uint _startICO, uint _startPostICO) public onlyOwner {
    require(now < start);
    require(
      _startPreICO > now &&
      _startICO > _startPreICO.add(28 days) &&
      _startPostICO > _startICO.add(28 days)
    );
    preICO.start = _startPreICO;
    ICO.start = _startICO;
    postICO.start = _startPostICO;
  }

  function setProxyMultisig(address _to) public onlyOwner returns(bool) {
    require(_to != address(0));
    proxyMultisig = _to;
    return true;
  }
    
}
