import WalletView from './WalletView.js';

export default class WalletController {

	static openMain(data ) {
		WalletView.renderMain(data);
		WalletController.getBalance();
	}

	static getBalance() {
		const data = {balance: 0};
		web3.eth.getAccounts(function (err, res) {
			if (err) {
				console.log(err);
			} else if (res.length < 1) {
				data.msg = "You are not authorized";
				WalletView.renderBalance(data);
			} else {
				data.address = res[0];
				web3.eth.getBalance(web3.eth.defaultAccount, function (err, res) {
					if (err) {
						console.log(err);
					} else {
						data.balance = web3.fromWei(res).toNumber();
						WalletView.renderBalance(data);
					}
				});
			}
		})
	}


	static async sendTx(data = {}) {
		ethereum.request({
			method: 'eth_sendTransaction',
			params: [{
				"to": data.to,
				"from": data.from,
				"value": web3.toHex(web3.toWei(data.value, 'ether')),
				"gas": data.gas.toString(16),
				"chainId": '4', //rinkeby
			}],
		}).then((txHash)=> {
			const link = `https://rinkeby.etherscan.io/tx/${txHash}\n`;
			WalletController.openMain({link: link});

		}, ()=>{
			console.log("Tx sending failed");
			WalletController.openMain({msg: "Failed!"});})
			.catch((error) => console.log("Something went wrong"));
	}
}
