/**
 * Copyright 2019 TIBCO Software Inc. All rights reserved.
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
 *
 * File name : ResultSetImpl.${EXT}
 * Created on: 05/01/2016
 * Created by: chung
 * SVN Id: $Id: ResultSetImpl.java 3149 2019-04-26 00:45:37Z sbangar $
 */

package com.tibco.tgdb.query.impl;

import java.util.ArrayList;
import java.util.Collection;
import java.util.Collections;
import java.util.List;

import com.tibco.tgdb.connection.TGConnection;
import com.tibco.tgdb.exception.TGException;
import com.tibco.tgdb.log.TGLogManager;
import com.tibco.tgdb.log.TGLogger;
import com.tibco.tgdb.query.TGResultSet;

public class ResultSetImpl<T> implements TGResultSet<T> {
    static TGLogger gLogger        = TGLogManager.getInstance().getLogger();

    private TGConnection conn;
    private int resultId;
    private List<T> resultList = new ArrayList<T>();
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

    public void addEntityToResultSet(T entity) {
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
	public T first() {
		if (isOpen == false) {
			return null;
		}
		if (resultList.size() == 0) {
			return null;
		}
		return resultList.get(0);
	}

	@Override
	public T last() {
		if (isOpen == false) {
			return null;
		}
		if (resultList.size() == 0) {
			return null;
		}
		return resultList.get(resultList.size() - 1);
	}

	@Override
	public T prev() {
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
	public T next() {
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
	public T getAt(int position) {
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
	
	Collection<T> getResultCollection() {
		return resultList;
	}
	
	public Collection<T> toCollection() {
		return Collections.unmodifiableCollection(resultList);
		
	}
}
