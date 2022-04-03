const dapr = require("dapr-client");

const DaprClient = dapr.DaprClient;
const HttpMethod = dapr.HttpMethod;
const CommunicationProtocolEnum = dapr.CommunicationProtocolEnum;

const mySidecarHost = "127.0.0.1";
const mySidecargRPCPort = process.env.DAPR_GRPC_PORT; // Note that the DAPR_GRPC_PORT environment variables is set by DAPR itself. https://docs.dapr.io/reference/environment/

const someServiceAppId = "some-service";

const client = new DaprClient(mySidecarHost, mySidecargRPCPort, CommunicationProtocolEnum.GRPC);

async function start() {
    const method = "echo"
    const r = await client.invoker.invoke(someServiceAppId, method, HttpMethod.POST, 
                                           { hello: "world" }
                                         );
    console.log(`after calling method ${method} on service ${someServiceAppId} - 
                 the response received was ${JSON.stringify(r)}`)
}

start().catch((e) => {
    console.error(e);
    process.exit(1);
});