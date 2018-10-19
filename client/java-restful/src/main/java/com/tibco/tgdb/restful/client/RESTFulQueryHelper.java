package com.tibco.tgdb.restful.client;

import java.io.BufferedReader;
import java.io.InputStreamReader;
import java.io.OutputStream;
import java.net.HttpURLConnection;
import java.net.URL;
import java.net.URLEncoder;
import java.util.List;
import java.util.Map;

import org.apache.http.HttpResponse;
import org.apache.http.client.HttpClient;
import org.apache.http.client.methods.HttpGet;
import org.apache.http.impl.client.DefaultHttpClient;
import org.json.JSONObject;

@SuppressWarnings("deprecation")
public abstract class RESTFulQueryHelper {
	
	public static String put(String url,  JSONObject json) throws Exception {
		System.out.println("xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx4");
		System.out.println(json.toString());
		System.out.println("xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx4");
        URL urlObj = new URL(url);
        HttpURLConnection conn = (HttpURLConnection) urlObj.openConnection();
        conn.setDoOutput(true);
        conn.setRequestMethod("PUT");
        conn.setRequestProperty("Content-Type", "application/json");
        OutputStream os = conn.getOutputStream();
        os.write(json.toString().getBytes());
        os.flush();
        if (conn.getResponseCode() != HttpURLConnection.HTTP_OK) {
            throw new RuntimeException("Failed : HTTP error code : "
                    + conn.getResponseCode());
        }
        BufferedReader br = new BufferedReader(new InputStreamReader(
                (conn.getInputStream())));
        
		 StringBuffer result = new StringBuffer();
		 String line = "";
		 while (null!=(line=br.readLine())) {
			 result.append(line);
		 }
		 System.out.println(result.toString());
        conn.disconnect();
		return result.toString();
	}
	
	public static String get(String url, Map<String, List<String>> requestPara) throws Exception {
		System.out.println("[RESTFulQueryHelper::get] requestPara = " + requestPara);
		
		 if(null!=requestPara) {
			 url = String.format("%s?%s", url, buildQueryParameter(requestPara));
			System.out.println("[RESTFulQueryHelper::get] url = " + url);
		 }
		 HttpGet request = new HttpGet(url);

		 HttpClient client = new DefaultHttpClient();
		 HttpResponse response = client.execute(request);
		 BufferedReader rd = new BufferedReader (new InputStreamReader(response.getEntity().getContent()));
		 StringBuffer content = new StringBuffer();
		 String line = null;
		 while (null!=(line=rd.readLine())) {
			 content.append(line);
		 }
		System.out.println("[RESTFulQueryHelper::get] content = " + content.toString());

		 return content.toString();
	 }
	   
	private static String buildQueryParameter(Map<String, List<String>> prop) {
		StringBuffer sb = new StringBuffer();
		for(String key : prop.keySet()) {
			for(String value : prop.get(key)) {
				sb.append(String.format("&%s=%s", key, URLEncoder.encode(value)));
			}
		}
		return sb.toString();
	}
}
