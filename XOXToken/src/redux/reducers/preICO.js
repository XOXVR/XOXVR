const initialState = {
  id: 0,
  totalSupply: 0,
  remainderTokens: 0,
  minInvestment: 0,
  sumInvestment: 0,
  end: 0,
  bonusLessTen: 0,
  bonusMoreTen: 0,
  start: 0
};

const preICO = (state = initialState, action) => {
  switch (action.type) {
    case 'STAGE_PRE_ICO': {
      const data = action.payload.map(num => num.toNumber());
      const [
        id,
        totalSupply,
        remainderTokens,
        minInvestment,
        sumInvestment,
        end,
        bonusLessTen,
        bonusMoreTen,
        start,
      ] = data;
      return {
        ...state, id, totalSupply, remainderTokens, minInvestment, sumInvestment, end, bonusLessTen, bonusMoreTen, start
      };
    }
    default: return state;
  }
};

export default preICO;
