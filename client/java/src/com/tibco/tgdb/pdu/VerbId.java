package com.tibco.tgdb.pdu;

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
 * <p/>
 * File name :VerbId
 * Created on: 12/26/14
 * Created by: suresh
 * <p/>
 * SVN Id: $Id: VerbId.java 2576 2018-10-17 02:36:19Z ssubrama $
 */
public enum VerbId {

    /**
     * Ping Message - Heart beats
     */
    PingMessage(0, com.tibco.tgdb.pdu.impl.PingMessage.class),

    /**
     * HandShake Request/Response protocol
     */
    HandShakeRequest (1, com.tibco.tgdb.pdu.impl.HandshakeRequest.class),
    HandShakeResponse(2, com.tibco.tgdb.pdu.impl.HandshakeResponse.class),

    /**
     * Authenticate Request/Response protocol
     */
    AuthenticateRequest(3, com.tibco.tgdb.pdu.impl.AuthenticateRequest.class),
    AuthenticateResponse(4, com.tibco.tgdb.pdu.impl.AuthenticateResponse.class),

    /**
     * Transaction - begin/commit/rollback protocol verbs
     */
    BeginTransactionRequest(5, com.tibco.tgdb.pdu.impl.BeginTransactionRequest.class),
    BeginTransactionResponse(6, com.tibco.tgdb.pdu.impl.BeginTransactionResponse.class),
    CommitTransactionRequest(7, com.tibco.tgdb.pdu.impl.CommitTransactionRequest.class),
    CommitTransactionResponse(8, com.tibco.tgdb.pdu.impl.CommitTransactionResponse.class),
    RollbackTransactionRequest(9, com.tibco.tgdb.pdu.impl.RollbackTransactionRequest.class),
    RollbackTransactionResponse(10, com.tibco.tgdb.pdu.impl.RollbackTransactionResponse.class),

    /**
     * Query Request/Response verbs
     */
    QueryRequest(11, com.tibco.tgdb.pdu.impl.QueryRequest.class),
    QueryResponse(12, com.tibco.tgdb.pdu.impl.QueryResponse.class),
    
    /**
     * Graph Traversal verbs
     */
    TraverseRequest(13, com.tibco.tgdb.pdu.impl.TraverseRequest.class),
    TraverseResponse(14, com.tibco.tgdb.pdu.impl.TraverseResponse.class),

    /**
     * Retrieve meta data
     */
    MetadataRequest(19, com.tibco.tgdb.pdu.impl.MetadataRequest.class),
    MetadataResponse(20, com.tibco.tgdb.pdu.impl.MetadataResponse.class),

    /**
     * Get entities
     */
    GetEntityRequest(21, com.tibco.tgdb.pdu.impl.GetEntityRequest.class),
    GetEntityResponse(22, com.tibco.tgdb.pdu.impl.GetEntityResponse.class),

    /**
     * Get LargeObject
     */
    GetLargeObjectRequest(23, com.tibco.tgdb.pdu.impl.GetLargeObjectRequest.class),
    GetLargeObjectResponse(24, com.tibco.tgdb.pdu.impl.GetLargeObjectResponse.class),

    /**
     * Import/Export verbs - They are admin request, and not supported by Java

    BeginExportRequest = 25,
    BeginExportResponse = 26,
    PartialExportRequest = 27,
    PartialExportResponse = 28,
    CancelExportRequest = 29,
    BeginImportRequest = 31,
    BeginImportResponse = 32,
    PartialImportRequest = 33,
    PartialImportResponse = 34,
     */

    /**
     * Dump Stacktrace request verb

    DumpStacktraceRequest = 39,
     */

    /**
     * Disconnect Request verbs
     */
    DisconnectChannelRequest(40, com.tibco.tgdb.pdu.impl.DisconnectChannelRequest.class),

    SessionForcefullyTerminated(41, com.tibco.tgdb.pdu.impl.SessionForcefullyTerminated.class),

    
    /**
     * Unknown Exception Message on the server.
     */
    ExceptionMessage(100, com.tibco.tgdb.pdu.impl.ExceptionMessage.class),


    InvalidMessage(-1, com.tibco.tgdb.pdu.impl.InvalidMessage.class);

    private Class<? extends TGMessage> msgClass;
    private short id;

    VerbId(int id, Class<? extends TGMessage> msgClass)
    {
        this.id = (short) id;
        this.msgClass = msgClass;
    }

    public static VerbId fromId(int id) {

        for (VerbId vid : VerbId.values()) {
            if (id == vid.id) return vid;
        }
        return InvalidMessage;
    }

    public Class<? extends TGMessage> getMessageClass()  {
        return msgClass;
    }

    public short getId() { return id; }
}
