/**
 * 
 */
var VerbId = require('../../../pdu/impl/VerbId').VerbId;

function testStub() {
	console.log('HANDSHAKE_REQUEST is ' + (VerbId.HANDSHAKE_REQUEST == VerbId.UTHENTICATE_REQUEST?'equals':'not equals') + ' to AUTHENTICATE_REQUEST')
	console.log('HANDSHAKE_REQUEST is ' + (VerbId.HANDSHAKE_REQUEST == VerbId.HANDSHAKE_REQUEST?'equals':'not equals') + ' to HANDSHAKE_REQUEST')
}

testStub();