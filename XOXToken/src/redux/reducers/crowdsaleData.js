const initialState = {
  multisig: '',
  decimals: 0,
  requireOnce: true,
  salesTokens: 0,
  rate: 0,
  softcap: 0,
  owner: '',
  addressCrowdsale: '',
  myWallet: '',
  deputyAddress: ''
};

const crowdsale = (state = initialState, action) => {
  switch (action.type) {
    case 'CROWDSALE': {
      const {
        multisig,
        decimals,
        requireOnce,
        salesTokens,
        rate,
        softcap,
        owner,
        addressCrowdsale,
        myWallet,
        deputyAddress
      } = action.payload;
      return {
        ...state, multisig, decimals, requireOnce, salesTokens, rate, softcap, owner, addressCrowdsale, myWallet, deputyAddress
      };
    }
    default: return state;
  }
};

export default crowdsale;
