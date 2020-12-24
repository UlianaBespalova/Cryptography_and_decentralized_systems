const _deploy_contracts = require("../migrations/2_deploy_contracts.js");

const NewToken = artifacts.require("NewToken");


contract("mintTest", accounts => {

    it ("mintTests: ", async()=> {

        let accTo = accounts[2];
        let instance = await NewToken.deployed();

        await instance.mint(accTo, 14);
        let toBalance = await instance.balanceOf(accTo);
        assert.equal(toBalance, 14);
    });
});


contract("burnTests", accounts => {

    it ("burnTest: ", async()=> {

        let accFrom = accounts[0];
        let instance = await NewToken.deployed();

        await instance.burn(accFrom, 4);
        let fromBalance = await instance.balanceOf(accFrom);
        assert.equal(fromBalance, 5000-4);
    });
});


contract("otherTests", accounts => {

    let instance;

    beforeEach(async function() {
        instance = await NewToken.deployed();
    });


    it ("transferTest: ", async()=> {
        let accFrom = accounts[0];
        let accTo = accounts[1];

        await instance.transfer(accTo, 7);

        let fromBalance = await instance.balanceOf(accFrom);
        let toBalance = await instance.balanceOf(accTo);

        assert.equal(fromBalance, 5000-7);
        assert.equal(toBalance, 7);
    });


    it ("mint&burnTest: ", async()=> { 

        let account = accounts[4];
        let err;

        try {
            await instance.mint(account, 1);
            await instance.burn(account, 1);
            
        } catch (error) {
            err = error;
        }        
        assert.notEqual(err, undefined, "Error must be thrown");
        assert.isAbove(err.message.search("can't call Mint and Burn methods in the same trx"), -1);
    });


    it ("mint&mintTest: ", async()=> { 

        let account = accounts[5];
        let err;

        try {
            await instance.mint(account, 1);
            await instance.mint(account, 1);
            
        } catch (error) {
            err = error;
        }        
        assert.equal(err, undefined);
    });
});