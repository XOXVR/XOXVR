pragma solidity ^0.4.23;

import './BurnableToken.sol';
import './StandardToken.sol';

contract XOXToken is BurnableToken, StandardToken {
  string public constant name = "XOX";
  string public constant symbol = "XOX";
  uint32 public constant decimals = 18;
  uint256 public constant INITIAL_SUPPLY = 1000E6 * (10 ** uint256(decimals));

  constructor(address _reserved) public {
    require(_reserved != address(0));

    totalSupply = INITIAL_SUPPLY;
    balances[msg.sender] = 900E6 * (10 ** uint256(decimals));
    balances[_reserved] = 100E6 * (10 ** uint256(decimals));
    emit Transfer(0x0, msg.sender, balances[msg.sender]);
    emit Transfer(0x0, _reserved, balances[_reserved]);
  }
}
