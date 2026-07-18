package com.qe.test.parser;

import java.util.Collections;
import java.util.LinkedHashMap;
import java.util.Map;
import java.util.regex.Pattern;

public class UrlParser {
    private static final Pattern URL_PATTERN = Pattern.compile(
            "^([a-zA-Z][a-zA-Z0-9+\\-.]*):" +                    // scheme
            "(?://(?:(.*?)@)?" +                                   // userinfo@
            "(\\[[a-fA-F0-9:]+\\]|[^:/?#]*)" +                    // host (IPv6 or regular)
            "(?::(\\d*))?" +                                       // port
            ")?(/[^?#]*)?" +                                       // path
            "(?:\\?([^#]*))?" +                                    // query
            "(?:#(.*))?$"                                          // fragment
    );

    private static final Pattern QUERY_PARAM_PATTERN = Pattern.compile("([^&=]+)=([^&]*)");

    private final String scheme;
    private final String userInfo;
    private final String host;
    private final int port;
    private final String path;
    private final String query;
    private final String fragment;
    private final Map<String, String> queryParams;
    private final String original;

    private UrlParser(Builder builder) {
        this.scheme = builder.scheme;
        this.userInfo = builder.userInfo;
        this.host = builder.host;
        this.port = builder.port;
        this.path = builder.path;
        this.query = builder.query;
        this.fragment = builder.fragment;
        this.queryParams = Collections.unmodifiableMap(builder.queryParams);
        this.original = builder.original;
    }

    public static UrlParser parse(String url) {
        if (url == null) {
            throw new NullPointerException("url must not be null");
        }
        var matcher = URL_PATTERN.matcher(url);
        if (!matcher.matches()) {
            throw new IllegalArgumentException("Malformed URL: " + url);
        }

        Builder builder = new Builder(url);
        builder.scheme = matcher.group(1).toLowerCase();
        builder.userInfo = matcher.group(2);

        String hostStr = matcher.group(3);
        if (hostStr != null) {
            if (hostStr.startsWith("[") && hostStr.endsWith("]")) {
                builder.host = hostStr.substring(1, hostStr.length() - 1);
            } else {
                builder.host = hostStr;
            }
        }

        String portStr = matcher.group(4);
        if (portStr != null && !portStr.isEmpty()) {
            try {
                builder.port = Integer.parseInt(portStr);
            } catch (NumberFormatException e) {
                throw new IllegalArgumentException("Invalid port: " + portStr, e);
            }
        } else {
            builder.port = getDefaultPort(builder.scheme);
        }

        String pathStr = matcher.group(5);
        builder.path = pathStr != null ? pathStr : "";

        builder.query = matcher.group(6);
        if (builder.query != null && !builder.query.isEmpty()) {
            builder.queryParams = parseQueryParams(builder.query);
        }

        builder.fragment = matcher.group(7);

        return builder.build();
    }

    private static Map<String, String> parseQueryParams(String query) {
        Map<String, String> params = new LinkedHashMap<>();
        var matcher = QUERY_PARAM_PATTERN.matcher(query);
        while (matcher.find()) {
            params.put(decodePercent(matcher.group(1)), decodePercent(matcher.group(2)));
        }
        return params;
    }

    static String decodePercent(String s) {
        try {
            return java.net.URLDecoder.decode(s, "UTF-8");
        } catch (java.io.UnsupportedEncodingException e) {
            return s;
        }
    }

    private static int getDefaultPort(String scheme) {
        return switch (scheme) {
            case "http" -> 80;
            case "https" -> 443;
            case "ftp" -> 21;
            case "ssh" -> 22;
            default -> -1;
        };
    }

    public String reconstruct() {
        StringBuilder sb = new StringBuilder();
        sb.append(scheme).append(":");
        if (host != null) {
            sb.append("//");
            if (userInfo != null) {
                sb.append(userInfo).append("@");
            }
            if (host.contains(":") && !host.startsWith("[")) {
                sb.append("[").append(host).append("]");
            } else {
                sb.append(host);
            }
            if (port != getDefaultPort(scheme) && port != -1) {
                sb.append(":").append(port);
            }
        }
        if (path != null && !path.isEmpty()) {
            sb.append(path);
        }
        if (query != null && !query.isEmpty()) {
            sb.append("?").append(query);
        }
        if (fragment != null && !fragment.isEmpty()) {
            sb.append("#").append(fragment);
        }
        return sb.toString();
    }

    public boolean isValid() {
        return scheme != null && !scheme.isEmpty();
    }

    public String getScheme() { return scheme; }
    public String getHost() { return host; }
    public int getPort() { return port; }
    public String getPath() { return path; }
    public String getQuery() { return query; }
    public String getFragment() { return fragment; }
    public String getUserInfo() { return userInfo; }
    public Map<String, String> getQueryParams() { return queryParams; }
    public String getOriginal() { return original; }

    private static class Builder {
        final String original;
        String scheme;
        String userInfo;
        String host;
        int port = -1;
        String path;
        String query;
        String fragment;
        Map<String, String> queryParams = Collections.emptyMap();

        Builder(String original) { this.original = original; }
        UrlParser build() { return new UrlParser(this); }
    }
}
