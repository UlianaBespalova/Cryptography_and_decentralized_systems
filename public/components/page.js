
export default class Page {

	static TxMessage(data={}) {
		return `<div> Success!</div>` +
			`<div>More info: <a style="display: inline" href=${data.link}>${data.link}</a></div>`
	}

	static sendTxForm() {
		return '<hr class="form_line">' +
			'<form action="#" class="form" name="sendTxForm">\n' +
			'            <div class="form__item">\n ' +
			'				<div class="err_msg" id="toError"></div>' +
			'                <label for="to" class="form__item__label">To</label>\n' +
			'                <input id="to" type="text" name="to" class="form__item__input">\n' +
			'            </div>\n' +
			'            <div class="form__item">\n ' +
			'				<div class="err_msg" id="valueError"></div>' +
			'                <label for="value" class="form__item__label">Value</label>\n' +
			'                <input id="value" type="text" name="value" class="form__item__input">\n' +
			'            </div>\n' +
			'            <div class="form__item">\n ' +
			'				<div class="err_msg" id="gasError"></div>' +
			'                <label for="gasLimit" class="form__item__label">Gas</label>\n' +
			'                <input id="gasLimit" type="text" value="21000" name="gasLimit" class="form__item__input">\n' +
			'            </div>\n' +
			'            <div class="form__button">\n' +
			'                <button id="sendButton" type="submit" name="send" class="page__sendTx_button page__sendTx_button-colorful">\n' +
			'                    Transfer\n' +
			'                </button>\n' +
			'            </div>\n' +
			'        </form>';
	}
}
