pragma solidity ^0.4.23;

import './SafeMath.sol';
import './XOXToken.sol';
import './Ownable.sol';

contract Crowdsale is Ownable {
  using SafeMath for uint;

  uint constant public SOFT_CAP = 1E24;
  uint constant private DECIMALS = 1E18;
  uint constant private PREICO_DURATION = 21 days;
  uint constant private ICO_DURATION = 28 days;

  uint public countRefund;
  uint public paginationNum = 200;
  uint public rate = 136900E9;
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
      _startICO > _startPreICO.add(PREICO_DURATION) &&
      _startPostICO > _startICO.add(ICO_DURATION)
    );
    require(requireOnce);

    requireOnce = false;
    multisig = owner;
    token = new XOXToken(_reserved);

    preICO = Stage({
      id: 0,
      totalSupply: 6E24,
      remainderTokens: 6E24,
      minInvestment: 1 ether,
      sumInvestment: 0,
      end: _startPreICO.add(PREICO_DURATION),
      bonusLessTen: 40,
      bonusMoreTen: 60,
      start: _startPreICO
    });

    ICO = Stage({
      id: 1,
      totalSupply: 54E24,
      remainderTokens: 54E24,
      minInvestment: 0.5 ether,
      sumInvestment: 0,
      end: _startICO.add(ICO_DURATION),
      bonusLessTen: 0,
      bonusMoreTen: 40,
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

    uint amountWei = msg.value;
    uint tokensValue = calculateTokens(amountWei);
    uint bonusPercent = amountWei < 10 ether ? stageICO.bonusLessTen : stageICO.bonusMoreTen;
    uint bonusTokens = tokensValue.mul(bonusPercent).div(100);
    uint tokensTotal = tokensValue.add(bonusTokens);

    if (tokensTotal <= stageICO.remainderTokens) {
      stageICO.remainderTokens = stageICO.remainderTokens.sub(tokensTotal);
    } else {
      tokensTotal = stageICO.remainderTokens;
      stageICO.remainderTokens = 0;
      uint paidTokens = tokensTotal.mul(uint(100).sub(bonusPercent)).div(100);
      amountWei = paidTokens.mul(rate).div(DECIMALS);
      uint remainderWei = msg.value.sub(amountWei);

      if (remainderWei > 0) {
        emit TransferWei(msg.sender, remainderWei);
        msg.sender.transfer(remainderWei);
      }
    }
    addStatisticsAndCalculateVal(tokensTotal, amountWei, msg.sender);

    if (salesTokens >= SOFT_CAP) {
      multisig.transfer(address(this).balance);
    }
  }

  function getValueForBackend(uint _btc, address _sender) public onlyOwnerAndProxyOwner returns(uint) {
    require(_btc > 0);
    require(_sender != address(0));
    require(now > ICO.end);

    uint tokensValue = calculateTokens(_btc);
    uint surrender;
    if (tokensValue <= postICO.remainderTokens) {
      postICO.remainderTokens = postICO.remainderTokens.sub(tokensValue);
    } else {
      tokensValue = postICO.remainderTokens;
      postICO.remainderTokens = 0;
      uint availableWei = tokensValue.mul(rate).div(DECIMALS);
      surrender = _btc.sub(availableWei);
    }
    token.transfer(_sender, tokensValue);
    return surrender;
  }

  function burnTokens() public onlyOwnerAndProxyOwner returns(bool) {
    require(ICO.end < now);

    uint remainderTokens = preICO.remainderTokens.add(ICO.remainderTokens);
    preICO.remainderTokens = 0;
    ICO.remainderTokens = 0;
    token.burn(this, remainderTokens);
    return true;
  }

  function refund() public onlyOwnerAndProxyOwner returns(bool) {
    require(ICO.end < now && salesTokens < SOFT_CAP);

    uint limit = countRefund.add(paginationNum);
    for(uint i = countRefund; i < investors.length; i++) {
      uint value = investments[investors[i]];
      investments[investors[i]] = 0;
      investors[i].transfer(value);
      if (i == limit) break;
    }
    countRefund = limit;
    return true;
  }

  function setDate(uint _startPreICO, uint _startICO, uint _startPostICO) public onlyOwner returns(bool) {
    require(now < preICO.start);
    require(
      _startPreICO > now &&
      _startICO > _startPreICO.add(PREICO_DURATION) &&
      _startPostICO > _startICO.add(ICO_DURATION)
    );

    preICO.start = _startPreICO;
    ICO.start = _startICO;
    postICO.start = _startPostICO;

    preICO.end =_startPreICO.add(PREICO_DURATION);
    ICO.end = _startICO.add(ICO_DURATION);
    return true;
  }

  function setMultisig(address _to) public onlyOwner returns(bool) {
    require(_to != address(0));

    multisig = _to;
    return true;
  }

  function setProxyMultisig(address _to) public onlyOwner returns(bool) {
    require(_to != address(0));

    proxyMultisig = _to;
    return true;
  }

  function setRate(uint _rate) public onlyOwnerAndProxyOwner returns(bool) {
    require(_rate > 0);

    rate = _rate.div(10);
    return true;
  }

  function getStage() internal view returns(Stage storage) {
    if (preICO.start <= now && now < preICO.end) {
      return preICO;
    } else if (ICO.start <= now && now < ICO.end) {
      return ICO;
    } else if (postICO.start <= now) {
      return postICO;
    } else {
      revert();
    }
  }

  function addStatisticsAndCalculateVal(uint _tokensVal, uint _ethVal, address _sender) internal {
    if (investments[_sender] == 0) investors.push(_sender);
    investments[_sender] = investments[_sender].add(_ethVal);
    salesTokens = salesTokens.add(_tokensVal);
    token.transfer(_sender, _tokensVal);
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
}
