const net = require("net");
const prompt = require('prompt-sync')()
// TODO: Allow the user to make packets and send them at will.


function createPacket() {
	var FullPacket = "";
	function GetID() {
		const ID = prompt("\nEnter a packet ID in decimal. ");
		if (isNaN(parseInt(ID)) || parseInt(ID) > 255) {
			console.log("\nInvalid ID. Try again");
			var p = GetID();
			return;
		}
		FullPacket += String.fromCharCode(parseInt(ID));
		return;
	}
	function GetData() {
		const Data = prompt("\nWrite a byte (a number from 0-255) for data. ");
		if (isNaN(parseInt(Data)) || parseInt(Data) > 255) {
			console.log("\nInvalid data input. Try again");
			var p = GetData()
			return p
		}
		FullPacket += String.fromCharCode(parseInt(Data)); 
		const again = prompt("\nAdd another byte? Type 'Yes' without quotes to do it.")
		if (again=="Yes") {
			GetData();
			return;
		}
		return;
	}
	GetID();
	GetData();
	var FinalPacket = String.fromCharCode((FullPacket.length + 1)) + FullPacket;
	return FinalPacket;
}



const client = net.createConnection({ port: 37669 }, () => {
	console.log("Oh golly");
	var packet = createPacket();
	var buffer = [];
	var stuff = new Buffer(packet, 'utf16le');
	for (var i = 0; i < stuff.length; i++) {
		buffer.push(stuff[i]);
	}
	console.log(buffer);
	console.log(packet);
	client.write(packet);
});

client.on('data', (data) => {
	console.log(data.toString());
	client.end()
});

client.on('end', () => {
	console.log("L");
});
