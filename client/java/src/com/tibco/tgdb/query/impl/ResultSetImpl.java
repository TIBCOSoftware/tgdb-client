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
 * File name : ResultSetImpl.${EXT}
 * Created on: 5/1/16
 * Created by: chung 
 * <p/>
 * SVN Id: $Id: ConnectionImpl.java 748 2016-04-25 17:10:38Z vchung $
 */

package com.tibco.tgdb.query.impl;

import java.util.ArrayList;
import java.util.List;

import com.tibco.tgdb.connection.TGConnection;
import com.tibco.tgdb.exception.TGException;
import com.tibco.tgdb.log.TGLogManager;
import com.tibco.tgdb.log.TGLogger;
import com.tibco.tgdb.model.TGEntity;
import com.tibco.tgdb.query.TGResultSet;


public class ResultSetImpl implements TGResultSet {
    static TGLogger gLogger        = TGLogManager.getInstance().getLogger();

    private TGConnection conn;
    private int resultId;
    private List<TGEntity> resultList = new ArrayList<TGEntity>();
    private boolean isOpen = false;
    private int currPos = 0;

    public ResultSetImpl(TGConnection conn, int resultId) {
        this.conn = conn;
        this.resultId = resultId;
        this.resultList = resultList;
        this.isOpen = true;
        this.currPos = -1;
    }

    public boolean isOpen() {
        return isOpen;
    }

    public void addEntityToResultSet(TGEntity entity) {
        resultList.add(entity);
    }

    public int getResultId() {
        return resultId;
    }

	@Override
	public boolean hasNext() {
		if (isOpen == false) {
			return false;
		}
		if (resultList.size() == 0) {
            return false;
        } else if (currPos < (resultList.size() - 1)) {
			return true;
		}
		return false;
	}

	@Override
	public void close() throws Exception {
		isOpen = false;
	}

	@Override
	public boolean hasExceptions() {
		// TODO Auto-generated method stub
		return false;
	}

	@Override
	public List<TGException> getExceptions() {
		// TODO Auto-generated method stub
		return null;
	}

	@Override
	public int count() {
		if (isOpen == false) {
			return 0;
		}
		return resultList.size();
	}

	@Override
	public TGEntity first() {
		if (isOpen == false) {
			return null;
		}
		if (resultList.size() == 0) {
			return null;
		}
		return resultList.get(0);
	}

	@Override
	public TGEntity last() {
		if (isOpen == false) {
			return null;
		}
		if (resultList.size() == 0) {
			return null;
		}
		return resultList.get(resultList.size() - 1);
	}

	@Override
	public TGEntity prev() {
		if (isOpen == false) {
			return null;
		}
		if (currPos > 0) {
			currPos--;
			return resultList.get(currPos);
		}
		return null;
	}

	@Override
	public TGEntity next() {
		if (isOpen == false) {
			return null;
		}
        if (resultList.size() == 0) {
            return null;
        } else if (currPos < (resultList.size() - 1)) {
			currPos++;
            return resultList.get(currPos);
		}
		return null;
	}

	@Override
	public int getPosition() {
		if (isOpen == false) {
			return 0;
		}
		return currPos;
	}

	@Override
	public TGEntity getAt(int position) {
		if (isOpen == false) {
			return null;
		}
		if (position >= 0 && position < resultList.size()) {
			return resultList.get(position);
		}
		return null;
	}

	@Override
	public void skip(int position) {
		if (isOpen == false) {
			return;
		}
		int newPos = currPos + position;
		if (newPos >=0 && newPos < resultList.size()) {
			currPos = newPos;
		}
	}
}
