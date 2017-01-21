/**
 * Copyright 2016 TIBCO Software Inc. All rights reserved.
 * 
 * Licensed under the Apache License, Version 2.0 (the "License"); You may not use this file except 
 * in compliance with the License.
 * A copy of the License is included in the distribution package with this file.
 * You also may obtain a copy of the License at http://www.apache.org/licenses/LICENSE-2.0
 * 
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

var ProtocolDataInputStream = require('../../pdu/impl/ProtocolDataInputStream').ProtocolDataInputStream,
    TGProtocolVersion       = require('../../TGProtocolVersion').TGProtocolVersion,
    TGException             = require('../../exception/TGException').TGException;
 
exports.VerbId = {
	
	PING_MESSAGE : { value : 0, name: 'PING_MESSAGE'},
	HANDSHAKE_REQUEST : { value : 1, name: 'HANDSHAKE_REQUEST'},
	HANDSHAKE_RESPONSE : { value : 2, name: 'HANDSHAKE_RESPONSE'},
	AUTHENTICATE_REQUEST : { value : 3, name: 'AUTHENTICATE_REQUEST'},
	AUTHENTICATE_RESPONSE : { value : 4, name: 'AUTHENTICATE_RESPONSE'},
	BEGIN_TRANS_REQUEST : { value : 5, name: 'BEGIN_TRANS_REQUEST'},
	BEGIN_TRANS_RESPONSE : { value : 6, name: 'BEGIN_TRANS_RESPONSE'},
	COMMIT_TRANS_REQUEST : { value : 7, name: 'COMMIT_TRANS_REQUEST'},
	COMMIT_TRANS_RESPONSE : { value : 8, name: 'COMMIT_TRANS_RESPONSE'},
	ROLLBCK_TRANS_REQUEST : { value : 9, name: 'ROLLBCK_TRANS_REQUEST'},
	ROLLBCK_TRANS_RESPONSE : { value : 10, name: 'ROLLBCK_TRANS_RESPONSE'},
	QUERY_REQUEST : { value : 11, name: 'QUERY_REQUEST'},
	QUERY_RESPONSE : { value : 12, name: 'QUERY_RESPONSE'},
	TRAVERSE_REQUEST : { value : 13, name: 'TRAVERSE_REQUEST'},
	TRAVERSE_RESPONSE : { value : 14, name: 'TRAVERSE_RESPONSE'},
	METADATA_REQUEST : { value : 19, name: 'METADATA_REQUEST'},
	METADATA_RESPONSE : { value : 20, name: 'METADATA_RESPONSE'},
	GET_ENTITY_REQUEST : { value : 21, name: 'GET_ENTITY_REQUEST'},
	GET_ENTITY_RESPONSE : { value : 22, name: 'GET_ENTITY_RESPONSE'},
	DISCONNECT_CHANNEL_REQUEST : { value : 23, name: 'DISCONNECT_CHANNEL_REQUEST'},
    
	EXCEPTION_MESSAGE : { value : 100, name: 'EXCEPTION_MESSAGE'},
	INVALID_MESSAGE : { value : -1, name: 'INVALID_MESSAGE'},
	verbIdFromBytes : function(buffer) {
		var inputStream = new ProtocolDataInputStream(buffer);
		var length = inputStream.readInt();

		if (length !== buffer.length) {
			throw new TGException('Buffer length mismatch.');
		}

		var magic = inputStream.readInt();
		if (magic !== TGProtocolVersion.getMagic()) {
			throw new TGException('Bad Message Magic');
		}

		var protocolVersion = inputStream.readShort();
		if (!TGProtocolVersion.isCompatible(protocolVersion)) {
			throw new TGException('Unsupported Protocol version');
		}

		return inputStream.readShort();
	}
};