/**
 * 
 */
var HexUtils                  = require('../../../utils/HexUtils').HexUtils,
	ProtocolDataOutputStream = require('../../../pdu/impl/ProtocolDataOutputStream').ProtocolDataOutputStream;

var position = 0;
var os = new ProtocolDataOutputStream();

function testStub01() {
	console.log("testStub01 entering .... ");
	os.writeByte(0);
	os.writeByte(1);
	console.log("After Write Byte : ");
	console.log(os._buffer);
	console.log(os.toBuffer());	
	os.writeBoolean(0);
	os.writeBoolean(1);
	console.log("After Write Boolean : ");
	console.log(os._buffer);
	console.log(os.toBuffer());
	os.writeShort(255);
	os.writeShort(300);
	console.log("After Write Shorts : ");
	console.log(os._buffer);
	console.log(os.toBuffer());
	os.writeInt(255);
	position = os.getPosition();
	os.writeInt(300);
	console.log("After Write Integers : ");
	console.log(os._buffer);
	console.log(os.toBuffer());	
	os.writeLong(0);
	os.writeLong(1);
	console.log("After Write Longs : ");
	console.log(os._buffer);
	console.log(os.toBuffer());	
}

function testUTF() {
	var os = new ProtocolDataOutputStream();
	os.writeUTF('Emily');
	console.log(HexUtils.formatHex(os._buffer));
	console.log(os._currentLength);	
}

function testFloat() {
	var os = new ProtocolDataOutputStream();
	os.writeFloat(3.141592);
	console.log(HexUtils.formatHex(os._buffer));
	console.log(os._currentLength);
}

function testLong() {
	var os = new ProtocolDataOutputStream();
	os.writeLong(Number.MAX_SAFE_INTEGER);
	os.writeLong(Number.MIN_SAFE_INTEGER);
	os.writeLong(-15);
	console.log(HexUtils.formatHex(os._buffer));
	console.log(os._currentLength);
}

function testDouble() {
	var os = new ProtocolDataOutputStream();
//	os.writeDouble(Number.MAX_Number);
	os.writeDouble(3.141592);
//	os.writeDouble(Number.MIN_Number);
	console.log(HexUtils.formatHex(os._buffer));
	console.log(os._currentLength);
}

function test() {
	testLong();
}

test();