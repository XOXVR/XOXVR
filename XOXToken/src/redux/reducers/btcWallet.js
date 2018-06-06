const initialState = {
  status: 'HERE WILL YOUR GENERATED ADDRESS'
};

const btcWallet = (state = initialState, action) => {
  let status;
  switch (action.type) {
    case 'SUCCESS':
      status = action.payload;
      return { ...state, status };
    case 'ERROR':
      status = action.payload;
      return { ...state, status };
    default: return initialState;
  }
};

export default btcWallet;
