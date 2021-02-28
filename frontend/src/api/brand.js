import axios from 'axios'
import BaseURL from './base'


const brandApi = {
    SearchAsYouType: async(term, cf) => {
        let url = `${BaseURL}/brand?term=${term}`
        console.log(url)
        axios.get(url)
            .then(function (response) {
                console.log(response.data)
                cf(response.data.data)
            })
            .catch(function (error) {
                console.log(error);
            })
    }
}

export default brandApi