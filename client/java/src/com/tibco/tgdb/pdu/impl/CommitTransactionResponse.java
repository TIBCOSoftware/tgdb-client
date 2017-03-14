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
 * File name : CommitTransactionResponse.java
 * Created on: 2/4/15
 * Created by: suresh 
 * <p/>
 * SVN Id: $Id: CommitTransactionResponse.java 1302 2017-01-13 22:31:31Z ssubrama $
 */


package com.tibco.tgdb.pdu.impl;

import com.tibco.tgdb.exception.TGException;
import com.tibco.tgdb.exception.TGTransactionException;
import com.tibco.tgdb.log.TGLogManager;
import com.tibco.tgdb.log.TGLogger;
import com.tibco.tgdb.log.TGLogger.TGLevel;
import com.tibco.tgdb.model.TGGraphObjectFactory;
import com.tibco.tgdb.pdu.TGInputStream;
import com.tibco.tgdb.pdu.TGOutputStream;
import com.tibco.tgdb.pdu.VerbId;

import java.io.IOException;
import java.util.ArrayList;
import java.util.List;

public class CommitTransactionResponse extends AbstractProtocolMessage {
    private static int lastStatus = 8000;
    public enum TransactionStatus {
        TGTransactionInvalid(-1),
        TGTransactionSuccess(0),
        TGTransactionAlreadyInProgress,
        TGTransactionClientDisconnected,
        TGTransactionMalFormed,
        TGTransactionGeneralError,
        TGTransactionVerificationError,
        TGTransactionInBadState,
        TGTransactionUniqueConstraintViolation,
        TGTransactionOptimisticLockFailed,
        TGTransactionResourceExceeded,
        TGCurrentThreadNotinTransaction,
        TGTransactionUniqueIndexKeyAttributeNullError;

        private int status;


        TransactionStatus() {
            this.status = ++CommitTransactionResponse.lastStatus; //Java stupidity
        }

        TransactionStatus(int status) {
            this.status = status;
        }

        public static TransactionStatus fromStatus(int status)
        {
            for (TransactionStatus ts : TransactionStatus.values()) {
                if (ts.status == status) return ts;
            }
            return TGTransactionInvalid;

        }
    }

    static TGLogger gLogger        = TGLogManager.getInstance().getLogger();

    private ArrayList<Long> addedIdList = null; 
    private ArrayList<Long> updatedIdList = null; 
    private ArrayList<Long> removedIdList = null; 
    private ArrayList<Integer> attrDescIdList = null; 
    private int attrDescCount = 0;
    private int addedCount = 0;
    private int updatedCount = 0;
    private int removedCount = 0;
    private TGGraphObjectFactory gof = null;
    private TGInputStream entityStream = null;
    private TGTransactionException exception;

    @Override
    protected void writePayload(TGOutputStream os) throws TGException, IOException {
    }

    @Override
    protected void readPayload(TGInputStream is) throws TGException, IOException {
    	gLogger.log(TGLevel.Debug, "Entering commit transaction response readPayload");
    	is.readInt(); // buf length
    	is.readInt(); // checksum
    	int status = is.readInt();// status code - currently zero
        exception = processTransactionStatus(is, status);
        if (exception != null) return;

        while (is.available() > 0) {
    		short opCode = is.readShort();
    		switch (opCode) {
    		case 0x1010: 
                attrDescIdList = new ArrayList<Integer>();
                attrDescCount = is.readInt();
                for (int i=0; i<attrDescCount; i++) {
                    int tempId = is.readInt();
                    attrDescIdList.add(tempId);
                    int realId = is.readInt();
                    attrDescIdList.add(realId);
                }
                gLogger.log(TGLevel.Debug, "Received %d attr desc", attrDescCount);
    			break;
    		case 0x1011:
                addedIdList = new ArrayList<Long>();
                addedCount = is.readInt();
                for (int i=0; i<addedCount; i++) {
                    Long tempId = is.readLong();
                    addedIdList.add(tempId);
                    Long realId = is.readLong();
                    addedIdList.add(realId);
                    long version = is.readLong();
                    addedIdList.add(version);
                }
                gLogger.log(TGLevel.Debug, "Received %d new entity", addedCount);
    			break;
            case 0x1012:
                updatedIdList = new ArrayList<Long>();
                updatedCount = is.readInt();
                for (int i=0; i<updatedCount; i++) {
                    Long id = is.readLong();
                    updatedIdList.add(id);
                    long version = is.readLong();
                    updatedIdList.add(version);
                }
                gLogger.log(TGLevel.Debug, "Received %d updated entity", updatedCount);
                break;
            case 0x1013:
                removedIdList = new ArrayList<Long>();
                removedCount = is.readInt();
                for (int i=0; i<removedCount; i++) {
                    Long id = is.readLong();
                    removedIdList.add(id);
                }
                gLogger.log(TGLevel.Debug, "Received %d delete results", removedCount);
                break;
            case 0x6789:
            	entityStream = is;
            	long pos = is.getPosition();
            	int count = is.readInt();
            	is.setPosition(pos);
                gLogger.log(TGLevel.Debug, "Received %d debug entities", count);
                return;
            default:
                break;
    		}
        }
    }

    private TGTransactionException processTransactionStatus(TGInputStream is, int status) throws TGTransactionException {
        TransactionStatus ts = TransactionStatus.fromStatus(status);
        String msg;
        switch (ts) {
            case TGTransactionSuccess:
                return null;

            case TGTransactionAlreadyInProgress:
            case TGTransactionClientDisconnected:
            case TGTransactionMalFormed:
            case TGTransactionGeneralError:
            case TGTransactionInBadState:
            case TGTransactionUniqueConstraintViolation:
            case TGTransactionOptimisticLockFailed:
            case TGTransactionResourceExceeded:
            case TGTransactionUniqueIndexKeyAttributeNullError:
            default:
                try {
                    msg = is.readUTF();
                }
                catch (Exception e) {
                    msg = "Error not available";
                }
                return TGTransactionException.buildException(ts, msg);
        }

    }

    @Override
    public boolean isUpdateable() {
        return false;
    }

    @Override
    public VerbId getVerbId() {
        return VerbId.CommitTransactionResponse;
    }

    public int getAttrDescCount() {
        return attrDescCount;
    }

    public int getAddedEntityCount() {
        return addedCount;
    } 

    public int getUpdatedEntityCount() {
        return updatedCount;
    } 

    public List<Integer> getAttrDescIdList() {
        return attrDescIdList;
    }

    public List<Long> getAddedIdList() {
        return addedIdList;
    }
    
    public List<Long> getUpdatedIdList() {
        return updatedIdList;
    }
    
    public TGInputStream getEntityStream() {
    	return entityStream;
    }

    public boolean hasException() { return exception != null;}

    public TGTransactionException getException() {return exception;}

}
