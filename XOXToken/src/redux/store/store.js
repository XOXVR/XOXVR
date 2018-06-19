import createSagaMiddleware from 'redux-saga';
import { createStore, applyMiddleware } from 'redux';
import { persistStore, persistReducer } from 'redux-persist';
import storage from 'redux-persist/lib/storage';

import reducers from '../reducers/combineReducers';
import sagas from '../saga/sagas';
import { middleware } from '../../router/Navigation';

const persistConfig = {
  key: 'root',
  storage,
  blacklist: ['statistics']
};

const sagaMiddleware = createSagaMiddleware();
const persistedReducer = persistReducer(persistConfig, reducers);

export const store = createStore(persistedReducer, applyMiddleware(sagaMiddleware, middleware));
export const persistor = persistStore(store);

sagaMiddleware.run(sagas);
