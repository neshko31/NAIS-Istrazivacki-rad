package com.example.stakeholders_service.choreography.event;

/**
 * Događaj koji auth servis objavljuje nakon što upiše korisnika u auth bazu.
 *
 * <p>U koreografisanom SAGA pattern-u, auth servis ne šalje KOMANDU stakeholder servisu
 * (kao u orkestrisanom), već objavljuje DOGAĐAJ koji opisuje šta se desilo.
 * Stakeholder servis SAM odlučuje da reaguje na ovaj događaj kreiranjem profila.
 *
 * <p>Routing key: "auth.user.created"
 * Exchange:      "saga.choreography.exchange"
 */
public class AuthUserCreatedEvent {

    private String sagaId;
    private Long   authId;
    private String username;
    private String email;
    private String role;
    private String firstname;
    private String lastname;

    public AuthUserCreatedEvent() {}

    public AuthUserCreatedEvent(String sagaId, Long authId, String username,
                                String email, String role,
                                String firstname, String lastname) {
        this.sagaId    = sagaId;
        this.authId    = authId;
        this.username  = username;
        this.email     = email;
        this.role      = role;
        this.firstname = firstname;
        this.lastname  = lastname;
    }

    public String getSagaId()             { return sagaId; }
    public void   setSagaId(String v)     { this.sagaId = v; }

    public Long   getAuthId()             { return authId; }
    public void   setAuthId(Long v)       { this.authId = v; }

    public String getUsername()           { return username; }
    public void   setUsername(String v)   { this.username = v; }

    public String getEmail()              { return email; }
    public void   setEmail(String v)      { this.email = v; }

    public String getRole()               { return role; }
    public void   setRole(String v)       { this.role = v; }

    public String getFirstname()          { return firstname; }
    public void   setFirstname(String v)  { this.firstname = v; }

    public String getLastname()           { return lastname; }
    public void   setLastname(String v)   { this.lastname = v; }
}
