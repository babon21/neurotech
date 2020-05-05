import {
  SET_ERRORS,
  GET_NEWS
} from '../types';
import axios from 'axios';


export const getNews = () => (dispatch) => {
  let myurl = 'http://localhost:8002/news'

  axios({
    method: 'get',
    url: myurl,
  }).then((res) => {
    if (res.data.error != null) {
      alert(res.data.error)
      dispatch({
        type: SET_ERRORS,
      })
      return
    }

    var titles = res.data.map(item => item.title)

    dispatch({
      type: GET_NEWS,
      payload: titles
    })
  }).catch((err) => {
    alert("getNews catch error", err)
    dispatch({
      type: SET_ERRORS,
    })
  });
};
