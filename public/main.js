'use strict';
import WalletController from './walletController.js';

const ethEnabled = () => {
	if (window.ethereum) {
		const defaultAccount = web3.eth.defaultAccount; //подключаемся к аккаунту MetaMask
		window.web3 = new Web3(window.ethereum);
		window.web3.eth.defaultAccount = defaultAccount;
		window.ethereum.enable();
		return true;
	}
	return false;
}

window.onload = () => {
	if (!ethEnabled()) {
		alert("You need to install MetaMask to use this service");
		return;
	}
	WalletController.openMain(); //после подключения к аккаунту делаем кнопку Send рабочей и запрашиваем текущий баланс
};
