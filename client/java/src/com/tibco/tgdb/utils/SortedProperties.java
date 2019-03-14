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

 * File name : SortedProperties.java
 * Created on: 2/5/15
 * Created by: suresh 
  * SVN Id: $Id: SortedProperties.java 2316 2018-04-26 23:49:37Z ssubrama $
 */


package com.tibco.tgdb.utils;


import java.util.Comparator;
import java.util.Map;
import java.util.SortedMap;
import java.util.TreeMap;

public class SortedProperties<K,V> extends TreeMap<K,V> implements TGProperties<K,V>
{

    public SortedProperties() {
        super();
    }

    public SortedProperties(Comparator<? super K> comparator) {
        super(comparator);
    }

    public SortedProperties(Map<? extends K, ? extends V> m) {
        super(m);
    }

    public SortedProperties(SortedMap<K, ? extends V> m) {
        super(m);
    }

    @Override
    public V getProperty(ConfigName cn) {
        V v = this.get(cn.getName());
        if (v == null) {
            v = this.get(cn.getAlias());
        }
        return v;
    }

    @Override
    public V getProperty(ConfigName cn, V defaultValue) {
        V v = getProperty(cn);
        return v == null ? defaultValue : v;
    }

    @Override
    public int getPropertyAsInt(ConfigName cn) {
        V v = getProperty(cn);
        if (v != null) return Integer.parseInt(v.toString());
        return 0;
    }

    @Override
    public int getPropertyAsInt(ConfigName cn, int defaultValue) {
        V v = getProperty(cn);
        if (v != null) return Integer.parseInt(v.toString());
        return defaultValue;
    }

    @Override
    public long getPropertyAsLong(ConfigName cn) {
        V v = getProperty(cn);
        if (v != null) return Long.parseLong(v.toString());
        return 0;
    }

    @Override
    public long getPropertyAsLong(ConfigName cn, long defaultValue) {
        V v = getProperty(cn);
        if (v != null) return Long.parseLong(v.toString());
        return defaultValue;
    }

    @Override
    public boolean getPropertyAsBoolean(ConfigName cn) {
        V v = getProperty(cn);
        if (v != null) return Boolean.parseBoolean(v.toString());
        return false;
    }

    @Override
    public boolean getPropertyAsBoolean(ConfigName cn, boolean defaultValue) {
        V v = getProperty(cn);
        if (v != null) return Boolean.parseBoolean(v.toString());
        return defaultValue;
    }
}
