const mongo = require('mongodb')
const axios = require('axios')
const chance = require('chance')
const Config = require("./config")

const chanceInstance = new chance()

async function sendToSearch(sBrands) {
    await axios
        .post(Config.BrandBulkInsertURL, sBrands)
        .then(res => {
            console.log(`statusCode: ${res.statusCode}`)
        })
        .catch(error => {
            console.error(error.response.data)
        })
}

function CtoSBrand(counter) {
    return {
        id: counter,
        version: 0,
        slug: chanceInstance.word({ length: 5 }),
        name:chanceInstance.word({ length: 5 }),
        description: chanceInstance.sentence({ words: 5 }),
        image_url: "https://picsum.photos/200/300",
    }
}

(async function main() {
    let counter = 1
    let data = []
    while (counter < 50000) {
        data = [...data, CtoSBrand(counter)]
        console.log(counter)
        if (counter % 500 == 0) {
            await sendToSearch(data)
            console.log(counter, " brands inserted")
            data = []
        }
       
        counter += 1
    }
    console.log("finished bulk inserting brands")
    process.exit(0)
})()

