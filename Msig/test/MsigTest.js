const _deploy_contracts = require("../migrations/2_deploy_contracts.js");

const Msig = artifacts.require("Msig");

contract("msigTests", accounts => {

    let instance;

    beforeEach(async function() {
        instance = await Msig.deployed();
    });


    it ("constructorTest: ", async()=> {

        threshold = await instance.m_threshold();
        num = await instance.m_numOwners();

        assert.equal(threshold.valueOf(), 1);
        assert.equal(num.valueOf(), 1);
    });


    it ("addOwnerTest: ", async()=> {
        
        let acc0 = accounts[0];
        let acc1 = accounts[1];
        let acc2 = accounts[2];

        await instance.addOwner(acc0); //попытаемся добавить адрес повторно
        num = await instance.m_numOwners();
        assert.equal(num.valueOf(), 1);

        await instance.addOwner(acc1);
        num = await instance.m_numOwners();
        assert.equal(num.valueOf(), 2);

        await instance.addOwner(acc2);
        num = await instance.m_numOwners();
        assert.equal(num.valueOf(), 3);
     });


     it ("removeOwnerTest: ", async()=> {
        
        let acc1 = accounts[1];
        let acc2 = accounts[2];

        await instance.removeOwner(acc1);
        num = await instance.m_numOwners();
        assert.equal(num.valueOf(), 2);

        await instance.removeOwner(acc1); //попытаемся удалить того же участника ещё раз
        num = await instance.m_numOwners();
        assert.equal(num.valueOf(), 2);

        await instance.removeOwner(acc2);
        num = await instance.m_numOwners();
        assert.equal(num.valueOf(), 1);
     });


     it ("changeThresholdTest: ", async()=> {
        
        let acc1 = accounts[1];

        await instance.changeThreshold(2); //превышает число участников, ничего не произойдёт 
        threshold = await instance.m_threshold();
        assert.equal(threshold.valueOf(), 1);

        await instance.addOwner(acc1);

        await instance.changeThreshold(2);
        threshold = await instance.m_threshold();
        assert.equal(threshold.valueOf(), 2);
     });


     it ("confirmAndCheckTest: ", async()=> {
        
        let op = "0x12345";

        res = await instance.confirmAndCheck.call(op);
        assert.equal(res, false);

        await instance.changeThreshold(1);
        res = await instance.confirmAndCheck.call(op);
        assert.equal(res, true);
     });
   });


   contract("transferTests", accounts => {

      let instance;
  
      beforeEach(async function() {
          instance = await Msig.deployed();
      });
  
  
      it ("depositTest: ", async()=> {
  
         total = await instance.getBalance();
         assert.equal(total.valueOf(), 0);

         await instance.depositEth({value: 15});

         total = await instance.getBalance();
         assert.equal(total.valueOf(), 15);

         await instance.depositEth({value: 12});

         total = await instance.getBalance();
         assert.equal(total.valueOf(), 27);
      });


      it ("depositErrorTest: ", async()=> {

         let acc = accounts[5];
         let err;
  
         oldTotal = await instance.getBalance();

         try {
            await instance.depositEth({from: acc, value: 20});
             
         } catch (error) {
             err = error;
         }        
         assert.notEqual(err, undefined, "Error must be thrown");

         newTotal = await instance.getBalance();
         assert.notEqual(oldTotal, newTotal);
      });


      it ("sendTest: ", async()=> {
  
         let acc = accounts[5];

         totalOld = await instance.getBalance();
         await instance.withdrawEth(acc, 3);
         totalNew = await instance.getBalance();

         assert.equal(totalOld.valueOf()-3, totalNew.valueOf());
      });


      it ("sendErrorTest: ", async()=> {

         let err;  
         oldTotal = await instance.getBalance();

         try {
            await instance.withdrawEth({value: 100500});
             
         } catch (error) {
             err = error;
         }        
         assert.notEqual(err, undefined, "Error must be thrown");

         newTotal = await instance.getBalance();
         assert.notEqual(oldTotal, newTotal);
      });
  

});


contract("msigTransferTests", accounts => {

   let instance;

   beforeEach(async function() {
       instance = await Msig.deployed();
   });


   it ("sendWithMsigTest: ", async()=> {

      let acc1 = accounts[1];
      let accNew = accounts[6];

      await instance.addOwner(acc1);
      await instance.changeThreshold(2);

      totalOld = await instance.getBalance();
      assert.equal(totalOld.valueOf(), 0);

      await instance.depositEth({value: 15});

      await instance.withdrawEth(accNew, 10);
      total1 = await instance.getBalance();
      assert.equal(totalOld.valueOf(), 0);
   });

});