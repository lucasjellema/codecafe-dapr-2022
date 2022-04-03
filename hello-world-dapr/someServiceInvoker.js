const dapr = require("dapr-client");

const DaprClient = dapr.DaprClient;
const HttpMethod = dapr.HttpMethod;
const CommunicationProtocolEnum = dapr.CommunicationProtocolEnum;

const mySidecarHost = "127.0.0.1";
const mySidecargRPCPort = process.env.DAPR_GRPC_PORT; // Note that the DAPR_GRPC_PORT environment variables is set by DAPR itself. https://docs.dapr.io/reference/environment/

const someServiceAppId = "some-service";

const client = new DaprClient(mySidecarHost, mySidecargRPCPort, CommunicationProtocolEnum.GRPC);

async function complexCalculation(a,b,c) {
    const method = "calculate"
    const r = await client.invoker.invoke(someServiceAppId, method, HttpMethod.POST, 
                                           { x: a, y:b, z:c }
                                         );
    console.log(`after calling method ${method} on service ${someServiceAppId} - 
                 the response received was ${JSON.stringify(r)}`)
    return r.outcome             
}

async function start() {
    const method = "echo"
    const r = await client.invoker.invoke(someServiceAppId, method, HttpMethod.POST, 
                                           { hello: "world" }
                                         );
    console.log(`after calling method ${method} on service ${someServiceAppId} - 
                 the response received was ${JSON.stringify(r)}`)
    const outcome = await complexCalculation(3.14159, 42, 2.71828)
    console.log(`the outcome of the complexCalculation ${outcome}`)                 
}

start().catch((e) => {
    console.error(e);
    process.exit(1);
});