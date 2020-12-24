const Msig = artifacts.require("Msig");

module.exports = function (deployer) {
  deployer.deploy(Msig, [], 1);
};