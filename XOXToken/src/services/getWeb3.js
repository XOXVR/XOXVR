import Web3 from 'web3';

const getweb3 = () => {
  let { web3 } = window;
  if (typeof web3 !== 'undefined') {
    web3 = new Web3(web3.currentProvider);
    console.log('Injected web3 detected.');
  } else {
    const provider = new Web3.providers.HttpProvider('http://localhost:9545');
    web3 = new Web3(provider);
    console.log('No web3 instance injected, using Local web3.');
  }
  return web3;
};

const web3Instance = getweb3();
export { web3Instance as default };
