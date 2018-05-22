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

  mapping(address => uint) public investments;
  
  XOXToken public token;

  struct Stage {
    uint totalSupply;
    uint remainderTokens;
    uint minInvestment;
    uint sumInvestment;
    uint duration;
    uint end;
    uint bonusLessTen;
    uint bonusMoreTen;
    uint id;
    uint start;
  }

  Stage public preICO;
  Stage public ICO;
  Stage public postICO;


  event TransferWei(address indexed addr, uint amount);

  constructor(uint _startPreICO, uint _startICO, uint _startPostICO, address _reserved) public {
    require(
      _startPreICO > now &&
      _startICO > _startPreICO.add(24 days) &&
      _startPostICO > _startICO.add(24 days)
    );
    require(requireOnce);

    requireOnce = false;
    multisig = owner;
    token = new XOXToken(_reserved);
    salesTokens = 0;
    countRefund = 0;
    start = _startPreICO;

    preICO = Stage({
      totalSupply: 6E24,
      remainderTokens: 6E24,
      minInvestment: 1 ether,
      sumInvestment: 0,
      duration: 24 days,
      end: start.add(24 days),
      bonusLessTen: 140,
      bonusMoreTen: 160,
      id: 0,
      start: _startPreICO
    });

    ICO = Stage({
      totalSupply: 54E24,
      remainderTokens: 54E24,
      minInvestment: 0.5 ether,
      sumInvestment: 0,
      duration: 24 days,
      end: _startICO.add(24 days),
      bonusLessTen: 0,
      bonusMoreTen: 140,
      id: 1,
      start: _startICO
    });

    postICO = Stage({
      totalSupply: 840E24,
      remainderTokens: 840E24,
      minInvestment: 0,
      sumInvestment: 0,
      duration: 0,
      end: 0,
      bonusLessTen: 0,
      bonusMoreTen: 0,
      id: 2,
      start: _startPostICO
    });
  }
  
  
  function() public payable {
    saleTokens();
  }
  

  function saleTokens() public payable {
    Stage storage stageICO = getStage();
    require(stageICO.id != 2);
    require(msg.value >= stageICO.minInvestment);
    
    uint ethVal = msg.value;
    address sender = msg.sender;
    uint tokensValue = calculateTokens(ethVal);
    
    if (tokensValue <= stageICO.remainderTokens) {
      stageICO.remainderTokens = stageICO.remainderTokens.sub(tokensValue);

      addStatisticsAndCalculateVal(tokensValue, ethVal, sender);
    } else {
        calculateTokensLessAvailable(stageICO, ethVal, sender);
    }
    withDrawal();
  }


  function getEthForBackend(uint _eth, address _sender) public onlyOwner {
    require(_eth > 0);
    require(_sender != address(0));
    
    uint ethVal = _eth;
    address sender = _sender;
    uint tokensValue = calculateTokens(ethVal);
    
    if (tokensValue <= postICO.remainderTokens) {
      postICO.remainderTokens = postICO.remainderTokens.sub(tokensValue);

      addStatisticsAndCalculateVal(tokensValue, ethVal, sender);
    } else {
        calculateTokensLessAvailable(postICO, ethVal, sender);
    }
    withDrawal();
  }
  
  
  function calculateTokensLessAvailable(Stage storage _stageNow, uint _ethValue, address _sender) internal {
      Stage storage stage = _stageNow;
      uint ethVal = _ethValue;
      address sender = _sender;
      uint tokensValue = stage.remainderTokens;
      stage.remainderTokens = 0;
      uint availableWei = tokensValue.mul(rate).div(DECIMALS);
      uint remainderWei = ethVal.sub(availableWei);

      addStatisticsAndCalculateVal(tokensValue, availableWei, sender);

      if (remainderWei > 0) {
        emit TransferWei(sender, remainderWei);
        sender.transfer(remainderWei);
      }
  }

  function addStatisticsAndCalculateVal(uint _tokensVal, uint _ethVal, address _sender) internal {
    uint ethVal = _ethVal;
    uint tokensVal = _tokensVal;
    address sender = _sender;
    
    if (investments[sender] == 0) investors.push(sender);
      investments[sender] = investments[sender].add(ethVal);
      salesTokens = salesTokens.add(tokensVal);
      token.transfer(sender, tokensVal);
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


  function calculateTokens(uint _ether) public view returns(uint) {
    return _ether.mul(DECIMALS).div(rate);
  }
    

  function calculateTotal(uint _tokens) public view returns (uint) {
    return _tokens.mul(rate);
  }
  
  
  function returnNumberInvestors() public view returns(uint) {
    return investors.length;
  }
    

  function refund() public onlyOwner {
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
    
}
