let a = 3.84911 * 1000;
let b = 6.24411 * 1000;
let c = a + b;
let d = 10708.219833 - c;
console.log(d);
// 54775372.00007

let i = 54775372.00007
let f = i - (2555 * 1000)
console.log(f);


const createNewOrder = async () => {
    try{
        let response = await fetch("http://localhost:8000/createOrder", {
            method:"POST",
            headers:{
                "content-type":"application/json"
            },
            body:JSON.stringify({
                "order_id":123,
                "token":"usdt",
                "amount":126,
                "callback":"url",
                "receiving_address":"TWYywngN3EfYiyY2NHzAHi4ad9B1uJNb8Y",
                "contract":"TY1DBj7Ys1bDcK37kwATaQpHxdTCnYrr1f",
                "receiving_private":"17c112793ba29f39dc0b6056695746a76f19bd8eb1e695d88d3c2dfdb30edb42",
                "callbackcount":1
            })
        })
        console.log("object");
    }catch(err){
        console.log("err", err);
    }
}

// for(let i=0; i<200; i++){
//     createNewOrder();
// }