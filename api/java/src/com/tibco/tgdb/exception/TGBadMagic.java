/**
 * Copyright 2019 TIBCO Software Inc.
 * All rights reserved.
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
 * <p/>
 * File name: TGBadMagic.java
 * Created on: 2014-12-27
 * Created by: suresh
 * <p/>
 * SVN Id: $Id: TGBadMagic.java 3996 2020-05-16 01:09:49Z vchung $
 */

package com.tibco.tgdb.exception;

public class TGBadMagic extends TGException {
    public TGBadMagic(String reason) {
        super(reason, TGExceptionType.BadMagic);
    }
}
