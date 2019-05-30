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
 *  File name : DatabaseStatisticsImpl.java
 *  Created on: 03/28/2019
 *  Created by: nimish
 *  
 *  
 *  SVN Id: $Id: DatabaseStatisticsImpl.java 3122 2019-04-25 21:38:58Z nimish $
 * 
 */

package com.tibco.tgdb.admin.impl;

import com.tibco.tgdb.admin.TGDatabaseStatistics;

public class DatabaseStatisticsImpl implements TGDatabaseStatistics {
	
	
	
	@Override
	public String toString() {
		return "DatabaseStatisticsImpl [dbSize=" + dbSize + ", numDataSegments=" + numDataSegments + ", dataSize="
				+ dataSize + ", dataUsed=" + dataUsed + ", dataFree=" + dataFree + ", dataBlockSize=" + dataBlockSize
				+ ", numIndexSegments=" + numIndexSegments + ", indexSize=" + indexSize + ", indexUsed=" + indexUsed
				+ ", indexFree=" + indexFree + ", blockSize=" + blockSize + "]";
	}
	protected long dbSize;
	protected int numDataSegments;
	protected long dataSize;
	protected long dataUsed;
	protected long dataFree;
	protected int dataBlockSize;
	
	protected int numIndexSegments;
	protected long indexSize;
	protected long indexUsed;
	protected long indexFree;
	protected int blockSize;
	
	public DatabaseStatisticsImpl(long _dbSize, int _numDataSegments, long _dataSize, long _dataUsed, long _dataFree,
			int _dataBlockSize, int _numIndexSegments, long _indexSize, long _indexUsed, long _indexFree, int _blockSize) {
		this.dbSize = _dbSize;
		this.numDataSegments = _numDataSegments;
		this.dataSize = _dataSize;
		this.dataUsed = _dataUsed;
		this.dataFree = _dataFree;
		this.dataBlockSize = _dataBlockSize;
		this.numIndexSegments = _numIndexSegments;
		this.indexSize = _indexSize;
		this.indexUsed = _indexUsed;
		this.indexFree = _indexFree;
		this.blockSize = _blockSize;
	}
	
	public long getDbSize() {
		return dbSize;
	}
	public int getNumDataSegments() {
		return numDataSegments;
	}
	public long getDataSize() {
		return dataSize;
	}
	public long getDataUsed() {
		return dataUsed;
	}
	public long getDataFree() {
		return dataFree;
	}
	public int getDataBlockSize() {
		return dataBlockSize;
	}
	public int getNumIndexSegments() {
		return numIndexSegments;
	}
	public long getIndexSize() {
		return indexSize;
	}
	public long getIndexUsed() {
		return indexUsed;
	}
	public long getIndexFree() {
		return indexFree;
	}
	public int getBlockSize() {
		return blockSize;
	}
	
	
	
	
}
