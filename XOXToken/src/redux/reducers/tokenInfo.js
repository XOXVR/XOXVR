const initialState = {
  name: '',
  symbol: '',
  decimals: 0,
  tokenAddress: '',
  initialSupply: 0,
  myTokens: 0
};

const token = (state = initialState, action) => {
  switch (action.type) {
    case 'TOKEN_INFO': {
      const {
        name,
        symbol,
        decimals,
        tokenAddress,
        initialSupply,
        myTokens
      } = action.payload;
      return {
        ...state, name, symbol, decimals, tokenAddress, initialSupply, myTokens
      };
    }
    default: return state;
  }
};

export default token;
