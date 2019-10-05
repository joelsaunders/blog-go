import axios from 'axios';

axios.defaults.xsrfCookieName = 'csrftoken';
axios.defaults.xsrfHeaderName = 'X-CSRFToken';

export default axios.create({
    baseURL: process.env.REACT_APP_API_URL
})