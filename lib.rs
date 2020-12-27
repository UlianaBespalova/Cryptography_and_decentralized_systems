#![cfg_attr(not(feature="std"), no_std)]

use ink_lang as ink;


#[ink::contract]
mod invoice {

    use std::time::SystemTime;


    #[ink(storage)]
    pub struct Invoice {
        invoice_amount: u128,
        paid_amount: u128,
        validity_period: u128,
        beneficiary: AccountId,
        payer: AccountId,
        partial_receiver: AccountId,
    }

    /// Indicates whether the invoice is already payed
    #[derive(Debug, PartialEq, scale::Encode)]
    #[cfg_attr(feature="std", derive(scale_info::TypeInfo))]
    pub enum Status { 
        Active, 
        Overdue, 
        Paid 
    }

    /// Errors that can occur upon calling this contract
    #[derive(Debug, PartialEq, scale::Encode)]
    #[cfg_attr(feature="std", derive(scale_info::TypeInfo))]
    pub enum Error {
        TransferFailed,
        InsufficientFunds,
        BelowSubsistenceThreshold,
    }


    /// Emitted when it is required to return the change to the payer
    #[ink(event)]
    pub struct Refund {
        #[ink(topic)]
        receiver: AccountId,
        #[ink(topic)]
        amount: u128,
    }

    /// Emitted when a payer send the amount (part of it) to the contract address
    #[ink(event)]
    pub struct Payment {
        #[ink(topic)]
        from: AccountId,
        #[ink(topic)]
        amount: u128,
    }

    /// Emitted when partial_receiver withdraws funds stored on the address of the contract
    #[ink(event)]
    pub struct Withdraw {
        #[ink(topic)]
        receiver: AccountId,
        #[ink(topic)]
        amount: u128,
    }





    impl Invoice {
        
        /// Creates a new Invoice contract with a specified participants (payer and beneficiary).
        ///
        /// # Panics
        /// 
        /// If validity period of the contract is less than the current moment.
        /// If specified partial receiver is not a participants of the contract.
        #[ink(constructor)]
        pub fn new( _invoice_amount: u128, _beneficiary: AccountId, _payer: AccountId,
                    _validity_period: u128, _partial_receiver: AccountId) -> Self {
            
            ensure_period_is_valid(&_validity_period);
            assert!(_partial_receiver == _payer || _partial_receiver == _beneficiary);
            Self {
                invoice_amount: _invoice_amount,
                beneficiary: _beneficiary,
                payer: _payer,
                validity_period: _validity_period,
                partial_receiver: _partial_receiver,
                paid_amount: 0,
            }            
        }


        /// Returnes the contract account balance.
        #[ink(message)]
        pub fn get_balance(&self) -> u128 {
            self.env().balance()
        }


        /// Returnes a status value of the invoice.
        #[ink(message)]
        pub fn get_status(&self) -> Status {
            if self.paid_amount == self.invoice_amount {
                Status::Paid
            } else if self.validity_period != 0 && self.validity_period < current_time() {
                Status::Overdue
            } else {
                Status::Active
            }
        }


        /// Transfers `amount` from the caller's account (payer) to the contract account.
        /// Returns the change if the payed amount exceeds the `paid_amount`
        /// 
        /// # Panics
        /// 
        /// If caller is not a payer.
        /// If the invoice status is not `Active`.
        #[ink(message, payable)]
        pub fn pay(&mut self, amount: u128) -> Result<(), Error>  {

            let from = self.env().caller();

            assert_eq! (from, self.payer);
            assert_eq! (self.get_status(), Status::Active);

            let will = self.paid_amount + amount;

            if will >= self.invoice_amount {
                if will > self.invoice_amount {
                    self.do_refund(will - self.invoice_amount);
                }
                self.paid_amount = self.invoice_amount;
                self.do_withdraw(self.beneficiary, self.get_balance());
            } else {
                self.paid_amount = will;
            }
            self.env().emit_event(Payment {
                from,
                amount,
            });  
            Ok(())
        }


        /// Allows `receiver` to withdraw `amount` from the contract's account.
        /// 
        /// # Panics
        /// 
        /// If there are not enough funds on the contract's account.
        /// If `partial_receiver` tries to withdraw funds from the paid contract.
        #[ink(message, payable)]
        pub fn withdraw(&mut self, receiver:AccountId, amount: u128) -> Result<(), Error>  {

            assert!(self.get_balance()>amount);

            let sender = self.env().caller();
            let status = self.get_status();

            assert!(status == Status::Paid && sender == self.beneficiary ||
                    status == Status::Active && self.validity_period == 0 && sender == self.partial_receiver ||
                    status == Status::Overdue && sender == self.partial_receiver);
            self.do_withdraw(receiver, amount)      
        }




        /// Transfers `amount` from the contract account to the caller's account (payer).
        /// 
        /// # Errors
        /// 
        /// Returns `InsufficientFunds` error if there are not enough funds on the contract's account.
        /// Returns `TransferFailed` error if the transfer of the funds failed.
        fn do_refund(&mut self, amount: u128) -> Result<(), Error> {
            
            let receiver = self.env().caller();

            if amount > self.env().balance() {
                return Err(Error::InsufficientFunds)
            }
            self.env().transfer(receiver, amount).map_err(|err| { 
                Error::TransferFailed 
            });
            self.env().emit_event(Refund {
                receiver,
                amount,
            });    
            Ok(())
        }


        /// Transfers `amount` from the contract account to the receiver's account.
        /// 
        /// # Errors
        /// 
        /// Returns `InsufficientFunds` error if there are not enough funds on the contract's account.
        /// Returns `TransferFailed` error if the transfer of the funds failed.
        fn do_withdraw(&mut self, receiver: AccountId, amount: u128) -> Result<(), Error> {

            if amount > self.env().balance() {
                return Err(Error::InsufficientFunds)
            }
            self.env().transfer(receiver, amount).map_err(|err| {
                Error::TransferFailed
            });
            self.env().emit_event(Withdraw {
                receiver,
                amount,
            });    
            Ok(())
        }
    }

/// Panic if validity period of the contract is less than the current moment.
fn ensure_period_is_valid (period: &u128) {
    let current_time = current_time();
    assert!(*period> current_time || *period == 0);
}

fn current_time() -> u128 {
    let now = SystemTime::now();
    let current_time = now.duration_since(std::time::UNIX_EPOCH).expect("Wrong time").subsec_millis() as u128;
    current_time/1000
}




    #[cfg(test)]
    mod tests {
        use super::*;
        use ink_lang as ink;
        use ink_env::test;
        use ink_env::call;

        type Accounts = test::DefaultAccounts<Environment>;
        fn default_accounts() -> Accounts {
            test::default_accounts()
                .expect("Test environment is expected to be initialized.")
        }

        fn set_sender(sender: AccountId) {
            const WALLET: [u8; 32] = [7; 32];
            test::push_execution_context::<Environment>(
                sender,
                WALLET.into(),
                1000000,
                1000000,
                test::CallData::new(call::Selector::new([0x00; 4])),
            );
        }
    
        #[ink::test]
        fn default_works() {
            let accounts = default_accounts();
            let invoice = Invoice::new(5, accounts.alice, accounts.bob, 10, accounts.alice);

            assert_eq!(invoice.invoice_amount, 5);
            assert_eq!(invoice.beneficiary, accounts.alice);
            assert_eq!(invoice.payer, accounts.bob);
            assert_eq!(invoice.validity_period, 10);
            assert_eq!(invoice.partial_receiver, accounts.alice);
            assert_eq!(invoice.paid_amount, 0);
        }


        #[ink::test]
        fn null_balance() {
            let accounts = default_accounts();
            let invoice = Invoice::new(5, accounts.alice, accounts.bob, 10, accounts.alice);
            let balance = invoice.get_balance();
            
            assert_eq!(balance, 0);
        }


        #[ink::test]
        #[should_panic]
        fn partial_receiver_fails() {
            let accounts = default_accounts();
            Invoice::new(5, accounts.alice, accounts.bob, 0, accounts.frank);
        }


        #[ink::test]
        fn return_active_status() {
            let accounts = default_accounts();
            let invoice = Invoice::new(5, accounts.alice, accounts.bob, 10, accounts.alice);
            let status = invoice.get_status();

            assert!(status == Status::Active);
        }


        #[ink::test]
        fn return_paid_status() {
            let accounts = default_accounts();
            let invoice = Invoice::new(0, accounts.alice, accounts.bob, 10, accounts.alice);
            let status = invoice.get_status();

            assert!(status == Status::Paid);
        }


        #[ink::test]
        #[should_panic]
        fn return_overdue_status_fails() {
            let accounts = default_accounts();
            let invoice = Invoice::new(0, accounts.alice, accounts.bob, 0, accounts.alice);
            let status = invoice.get_status();

            assert!(status == Status::Overdue);
        }


        #[ink::test]
        #[should_panic]
        fn pay_fails_not_payer() {
            let accounts = default_accounts();
            let mut invoice = Invoice::new(0, accounts.alice, accounts.bob, 0, accounts.alice);
            invoice.pay(12);
        }


        #[ink::test]
        #[should_panic]
        fn pay_fails_already_paid() {
            let accounts = default_accounts();
            let mut invoice = Invoice::new(0, accounts.alice, accounts.bob, 0, accounts.alice);
            invoice.pay(12);
        }


        #[ink::test]
        fn payed_part() {
            let accounts = default_accounts();
            let mut invoice = Invoice::new(100, accounts.alice, accounts.bob, 0, accounts.alice);

            set_sender(accounts.bob);
            let res = invoice.pay(12);
            assert_eq!(res, Ok(()));
            let balance = invoice.paid_amount;
            assert_eq!(balance, 12);
            assert_eq!(invoice.get_status(), Status::Active);
        }


        #[ink::test]
        fn payed_enough() {
            let accounts = default_accounts();
            let mut invoice = Invoice::new(100, accounts.alice, accounts.bob, 0, accounts.alice);

            set_sender(accounts.bob);
            let res = invoice.pay(120);
            assert_eq!(res, Ok(()));
            let balance = invoice.paid_amount;
            assert_eq!(balance, 100);
            assert_eq!(invoice.get_status(), Status::Paid);
        }


        #[ink::test]
        #[should_panic]
        fn withdraw_failes_no_balance() {
            let accounts = default_accounts();
            let mut invoice = Invoice::new(0, accounts.alice, accounts.bob, 0, accounts.alice);
            let res = invoice.withdraw(accounts.alice, 100);
        }


        #[ink::test]
        fn do_refund_failes_insufficient_funds() {
            let accounts = default_accounts();
            let mut invoice = Invoice::new(0, accounts.alice, accounts.bob, 0, accounts.alice);
            
            let res = invoice.do_refund(12);
            assert_eq!(res, Err(Error::InsufficientFunds));
        }


        #[ink::test]
        fn do_withdraw_failes_insufficient_funds() {
            let accounts = default_accounts();
            let mut invoice = Invoice::new(0, accounts.alice, accounts.bob, 0, accounts.alice);
            
            let res = invoice.do_withdraw(accounts.alice, 12);
            assert_eq!(res, Err(Error::InsufficientFunds));
        }

    }
}