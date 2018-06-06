export const generateWallet = data => ({ type: 'NEW_WALLET', payload: data });
export const isError = data => ({ type: 'ERROR', payload: data });
export const success = data => ({ type: 'SUCCESS', payload: data });
