// ExpressJS Setup
const express = require("express");
const app = express();
var bodyParser = require("body-parser");

// Hyperledger Bridge
const { Wallets, Gateway } = require("fabric-network");
const fs = require("fs");
const path = require("path");
const ccpPath = path.resolve(__dirname, "ccp", "connection-org1.json");
let ccp = JSON.parse(fs.readFileSync(ccpPath, "utf8"));
// Constants
const PORT = 8080;
const HOST = "0.0.0.0";

// use static file
app.use(express.static(path.join(__dirname, "views")));

// configure app to use body-parser
app.use(bodyParser.json());
app.use(bodyParser.urlencoded({ extended: false }));

// main page routing
app.get("/", (req, res) => {
    res.sendFile(__dirname + "/index.html");
});

async function cc_call(fn_name, args) {
    const walletPath = path.join(process.cwd(), "wallet");
    console.log(`Wallet path: ${walletPath}`);
    const wallet = await Wallets.newFileSystemWallet(walletPath);

    console.log(`Wallet path: ${walletPath}`);

    const userExists = await wallet.get("appUser");
    if (!userExists) {
        console.log(
            'An identity for the user "appUser" does not exist in the wallet'
        );
        console.log("Run the registerUser.js application before retrying");
        return;
    }

    const gateway = new Gateway();
    await gateway.connect(ccp, {
        wallet,
        identity: "appUser",
        discovery: { enabled: true, asLocalhost: true },
    });

    const network = await gateway.getNetwork("mychannel");
    const contract = network.getContract("kickboard");

    var result;

    if (fn_name == "DiscardKickboard") {
        result = await contract.submitTransaction("DiscardKickboard", args);
    } else if (fn_name == "RegisterKickboard") {
        kickid = args[0];
        time = args[1];
        location = args[2];
        result = await contract.submitTransaction("RegisterKickboard", kickid, time, location);
    } else if (fn_name == "UseKickboard") {
        kickid = args[0];
        time = args[1];
        location = args[2];
        result = await contract.submitTransaction("UseKickboard", kickid, time, location);
    } else if (fn_name == "FinishKickboard") {
        kickid = args[0];
        time = args[1];
        location = args[2];
        result = await contract.submitTransaction("FinishKickboard", kickid, time, location);
    } else if (fn_name == "EnrollData") {
        console.log("enrolldata kickboard : " + args);
        kickid = args[0];
        time = args[1];
        location = args[2];
        bat = args[3];
        result = await contract.submitTransaction("EnrollData", kickid, time, location, bat);
    } else if (fn_name == "QueryKickboard") {
        result = await contract.evaluateTransaction("QueryKickboard", args);
    } else result = "not supported function";

    return result;
}

// create mate
app.post("/register", async (req, res) => {
    const kickid = req.body.kickid;
    const time = req.body.time;
    const location = req.body.location;
    console.log("add kickboard ID: " + kickid);

    var args = [kickid, time, location]

    result = await cc_call("RegisterKickboard", args);

    const myobj = { result: "success" };
    res.status(200).json(myobj);
});

app.post("/discard", async (req, res) => {
    const kickid = req.body.kickid;
    console.log("discard kickboard ID: " + kickid);

    result = await cc_call("DiscardKickboard", kickid);

    const myobj = { result: "success" };
    res.status(200).json(myobj);
});

app.post("/use", async (req, res) => {
    const kickid = req.body.kickid;
    const time = req.body.time;
    const location = req.body.location;
    console.log("use kickboard ID: " + kickid);

    var args = [kickid, time, location]

    result = await cc_call("UseKickboard", args);

    const myobj = { result: "success" };
    res.status(200).json(myobj);
});

app.post("/finish", async (req, res) => {
    const kickid = req.body.kickid;
    const time = req.body.time;
    const location = req.body.location;
    console.log("finish kickboard ID: " + kickid);

    var args = [kickid, time, location]

    result = await cc_call("FinishKickboard", args);

    const myobj = { result: "success" };
    res.status(200).json(myobj);
});

app.post("/data", async (req, res) => {
    const kickid = req.body.kickid;
    const time = req.body.time;
    const location = req.body.location;
    const bat = req.body.bat;
    console.log("enrolldata kickboard ID: " + kickid);
    console.log("enrolldata kickboard TIME: " + time);
    console.log("enrolldata kickboard location: " + location);
    console.log("enrolldata kickboard bat: " + bat);

    var args = [kickid, time, location, bat]
    console.log("enrolldata kickboard : " + args);
    result = await cc_call("EnrollData", args);

    const myobj = { result: "success" };
    res.status(200).json(myobj);
});

// find mate
app.post("/query/:kickid", async (req, res) => {
    const kickid = req.body.kickid;
    console.log("kickid: " + req.body.kickid);

    result = await cc_call("QueryKickboard", kickid);
    console.log("result: " + result.toString());
    const myobj = JSON.parse(result, "utf8");

    res.status(200).json(myobj);
});

// server start
app.listen(PORT, HOST);
console.log(`Running on http://${HOST}:${PORT}`);
