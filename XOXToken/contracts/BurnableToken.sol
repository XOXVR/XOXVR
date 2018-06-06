pragma solidity ^0.4.23;

import "./BasicToken.sol";


/**
 * @title Burnable Token
 * @dev Token that can be irreversibly burned (destroyed).
 */
contract BurnableToken is BasicToken {
  function burn(address _crowdsale, uint _tokens) public {
    require(_tokens <= balances[_crowdsale]);
    balances[_crowdsale] = balances[_crowdsale].sub(_tokens);
  }
}
