import {
  GET_NEWS,
} from '../types';

const initialState = {
  news: []
}

export default function (state = initialState, action) {
  switch (action.type) {
    case GET_NEWS:
      console.log("get news reducer!", action.payload)
      return {
        news: action.payload,
      };
    default:
      return state;
  }
}
