// SPDX-License-Identifier: MIT
pragma solidity ^0.4.24;

import 'zeppelin-solidity/contracts/token/ERC20/StandardToken.sol';

contract NewToken is StandardToken {

    string public name = 'NewToken';
    string public symbol = 'NT';
    uint public decimals = 10;
    uint public INITIAL_SUPPLY = 5000;

    mapping(address => uint) internal minted;
    mapping(address => uint) internal burned;


    constructor() {
        totalSupply_ = INITIAL_SUPPLY;
        balances[msg.sender] = INITIAL_SUPPLY;
    }

    function mint (address account, uint amount) public returns (bool) {
        require(account != address(0), "Error: you can't mint to the contract address");

        require (burned[account] < now, "Error: you can't call Mint and Burn methods in the same trx");
        delete burned[account];
        minted[account] = now;

        totalSupply_ = totalSupply_.add(amount);
        balances[account] = balances[account].add(amount);

        emit Transfer(address(0), account, amount);
        return true;
    }


    function burn (address account, uint amount) public returns (bool) {
        require(account != address(0), "Error: you can't burn from the contract address");

        require (minted[account] < now, "Error: you can't call Mint and Burn methods in the same trx");
        delete minted[account];
        burned[account] = now;

        totalSupply_ = totalSupply_.sub(amount);
        balances[account] = balances[account].sub(amount);

        emit Transfer(account, address(0), amount);
        return true;
    }
}