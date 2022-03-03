

export const wsConnection = (url: string) => {
    const ws = new WebSocket(`ws:${url}`)
    ws.onopen = () => {
        console.log("web socket Connection established");
        ws.send(JSON.stringify({
            content: "howdy boys"
        }))
    }
    ws.onclose = (event) => {
        console.log("Socket closed connection", event);

    }
    ws.onerror = (error) => {
        console.log("Socket Error: ", error);
    }
    ws.onmessage = (message) => {
        console.log(message.data);

    }
    return (message)=>{
        if(message) ws.send(JSON.stringify(message))
    };
}