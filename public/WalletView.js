import page from './components/page.js'
import WalletController from "./walletController.js";

export default class WalletView {

	static renderMain(data) {
		if (data!==undefined) { //показать сообщение о прерыдущей операции, если она была
			WalletView._showTxMsg(data);
			const sendTxFormEl = document.getElementById("sendTxForm");
			if (sendTxFormEl) {
				sendTxFormEl.innerHTML = "";
			}
		}
		let sendTxFormOpened = false;
		const sendTxButtonEl = document.getElementById("sendTxButton"); //обработчик кнопки Send для открытия формы
		const sendTxFormEl = document.getElementById("sendTxForm");
		if (sendTxButtonEl && sendTxFormEl) {
			sendTxButtonEl.addEventListener('click', (event) => {
				event.preventDefault();
				if (this.accountData===undefined) {
					alert("The service is currently trying to connect to MetaMask. Please wait a little");
					return;
				}
				if (sendTxFormOpened) {
					sendTxFormEl.innerHTML = "";
				} else {
					sendTxFormEl.innerHTML = page.sendTxForm();
					this._addTxSubmitListener();
					const TxMsgEl = document.getElementById("TxInfo");
					if (TxMsgEl) {
						TxMsgEl.innerText = "";
					}
				}
				sendTxFormOpened = !sendTxFormOpened;
			})
		}
	}

	static renderBalance(data= null) {
		this.accountData = data;
		const accountInfoEl = document.getElementById("accountInfo");
		if (accountInfoEl && data) {
			accountInfoEl.innerText = `address: ${data.address}\n` +
				`balance: ${data.balance}`}
	}

	static _addTxSubmitListener() { //обработчик кнопки Transfer для отправки транзакции
		const toInput = document.forms.sendTxForm.to;
		const valueInput = document.forms.sendTxForm.value;
		const gasInput = document.forms.sendTxForm.gasLimit;

		document.forms.sendTxForm.addEventListener('submit', (event) => {
			event.preventDefault();
			const errMsgList = document.getElementsByClassName('err_msg');
			for (let i = 0; i < errMsgList.length; i += 1) {
				errMsgList[i].innerText = "";
			}

			//-------валидация
			const invalidInput = [];
			if (toInput.value.length!==42 || !toInput.value.startsWith("0x"))
				invalidInput.push({name: "toError", msg: "Invalid address"});

			const valueInputNum = Number.parseFloat(valueInput.value);
			if (isNaN(valueInputNum)||valueInputNum<0)
				invalidInput.push({name: "valueError", msg: "Invalid value"});
			if (valueInputNum>this.accountData.balance)
				invalidInput.push({name: "valueError", msg: "Not enough ETH to send"});

			const gasInputNum = Number.parseInt(gasInput.value, 10);
			if (isNaN(gasInputNum)||gasInputNum<0)
				invalidInput.push({name: "gasError", msg: "Invalid gas limit"});

			if (invalidInput.length>0) {
				this._showErrMsg(invalidInput);
				return
			}

			WalletController.sendTx({
				to: toInput.value,
				from: this.accountData.address,
				value: valueInputNum,
				gas: gasInputNum,
			});
		}, false);
	}

	static _showErrMsg(data=[]) {
		data.forEach((item, i, arr)=>{
			const msgEl = document.getElementById(item.name);
			if (msgEl)
				msgEl.innerText = item.msg;
		})
	}

	static _showTxMsg(data) {
		const TxMsgEl = document.getElementById("TxInfo");
		if (TxMsgEl) {
			TxMsgEl.innerHTML = page.TxMessage(data);
		}
	}
}



