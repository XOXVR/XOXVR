const initialState = {
  objInvest: 0
};

const statistics = (state = initialState, action) => {
  switch (action.type) {
    case 'STATISTICS_ADD': {
      const objInvest = action.payload;
      return {
        ...state, objInvest
      };
    }
    default: return state;
  }
};

export default statistics;
