// const mongo = require('mongodb')
// const axios = require('axios')
// const Config = require("./config")


// async function sendToSearch(sShops) {
//     await axios
//     .post(Config.ShopBulkInsertURL, sShops)
//     .then(res => {
//       console.log(`statusCode: ${res.statusCode}`)
//     })
//     .catch(error => {
//       console.error(error.response.data)
//     })
// }

// function CtoSShop(cShop) {
//     return {
//         id : cShop.id,
//         version: cShop.version,
//         slug: cShop.slug,
//         shop_name: cShop.name,
//         shop_image: cShop.image,
//         owner_number: cShop.owner_name,
//         owner_name: cShop.description,
//         contact_number: cShop.contact_number,
//     }
// }

// (async function main() {
//     let counter = 1
//     let data = []
//     while (counter < 50000) {
//         data = [...data, CtoSShop(counter)]
//         console.log(counter)
//         if (counter % 500 == 0) {
//             await sendToSearch(data)
//             console.log(counter, " shops inserted")
//             data = []
//         }
       
//         counter += 1
//     }
//     console.log("finished bulk inserting shops")
//     process.exit(0)
// })()

