package com.example.stakeholders_service.Utils;

import io.jsonwebtoken.Claims;
import io.jsonwebtoken.Jwts;
import io.jsonwebtoken.security.Keys;
import org.springframework.stereotype.Component;
import java.security.Key;

@Component
public class JwtUtil {

    // secret key
    private static final String SECRET = "your-secure-secret-key-which-is-so-lovely";
    // issuer
    private static final String EXPECTED_ISSUER = "AuthService";
    private Key getSigningKey() {
        return Keys.hmacShaKeyFor(SECRET.getBytes());
    }
    public Claims extractClaims(String token) {
        return Jwts.parserBuilder()
                .setSigningKey(getSigningKey())
                .requireIssuer(EXPECTED_ISSUER)
                .build()
                .parseClaimsJws(token)
                .getBody();
    }
    public Long extractUserId(String token) {
        Claims claims = extractClaims(token);
        return Long.valueOf(claims.get("user_id").toString());
    }

    public boolean isTokenValid(String token) {
        try {
            extractClaims(token); // baca exception ako nije validan
            return true;
        } catch (Exception e) {
            return false;
        }
    }
}