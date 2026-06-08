package com.example.stakeholders_service.Saga.command;

/**
 * Komanda koju prima stakeholder servis od SagaOrchestratora (auth servis).
 * Nalaže stakeholder servisu da kreira profil korisnika u svojoj bazi.
 *
 * Polja:
 *   sagaId    – korelacioni ID sage (vraća se u reply-u)
 *   authId    – ID korisnika iz auth baze (mapira se na registry_id u stakeholder bazi)
 *   username  – korisničko ime
 *   email     – email adresa
 *   role      – uloga (TOURIST, GUIDE, ADMIN)
 *   firstname – ime
 *   lastname  – prezime
 */
public class RegisterUserCommand {

    private String sagaId;
    private Long   authId;
    private String username;
    private String email;
    private String role;
    private String firstname;
    private String lastname;

    public RegisterUserCommand() {}

    public RegisterUserCommand(String sagaId, Long authId, String username,
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

    public String getSagaId()              { return sagaId; }
    public void   setSagaId(String v)      { this.sagaId = v; }

    public Long   getAuthId()              { return authId; }
    public void   setAuthId(Long v)        { this.authId = v; }

    public String getUsername()            { return username; }
    public void   setUsername(String v)    { this.username = v; }

    public String getEmail()               { return email; }
    public void   setEmail(String v)       { this.email = v; }

    public String getRole()                { return role; }
    public void   setRole(String v)        { this.role = v; }

    public String getFirstname()           { return firstname; }
    public void   setFirstname(String v)   { this.firstname = v; }

    public String getLastname()            { return lastname; }
    public void   setLastname(String v)    { this.lastname = v; }
}
