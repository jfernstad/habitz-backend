function retrieveTodaysHabitz(){
    // var request = new XMLHttpRequest();
    // var url = "/v1/today";

    // request.onreadystatechange = function () {
    //     if (this.readyState == 4 && this.status == 200) {
    //         console.log(this.responseText);
    //     }
    // };
    // request.open("POST", url);
    // request.setRequestHeader("Content-Type", "application/json;charset=UTF-8");
    // request.send(JSON.stringify(body));

    callback = function (status, payload){
        console.log("RESPONSE ["+status+"]: " + payload)
    }

    GET("/v1/today", callback)
}

function POST(url, payload, callback, errCallback = null){
    var request = new XMLHttpRequest();
    request.withCredentials = true;
    request.open("POST", url);
    request.setRequestHeader("Content-Type", "application/json;charset=UTF-8");
    request.setRequestHeader("Authorization", "Bearer " + localStorage.getItem("habitz-token"))

    request.onreadystatechange = function () {
        if (this.readyState == 4){
            if (this.status >= 200 && this.status <= 299) {
                callback(this.status, this.responseText)
            }
            else if (this.status >= 400 && errCallback != null) {
                errCallback(this.status, this.responseText)
            }
    }
    };
    request.send(JSON.stringify(payload));
}

function GET(url, callback, errCallback = null) {
    var request = new XMLHttpRequest();
    request.withCredentials = true;
    request.open("GET", url);
    request.setRequestHeader("Authorization", "Bearer " + localStorage.getItem("habitz-token"))

    request.onreadystatechange = function () {
        if (this.readyState == 4) {
            if (this.status >= 200 && this.status <= 299) {
                callback(this.status, this.responseText)
            }
            else if (this.status >= 400 && errCallback != null) {
                errCallback(this.status, this.responseText)
            }
        }
    };
    request.send();
}