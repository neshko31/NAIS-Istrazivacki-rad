package com.example.stakeholders_service.choreography.event;

/**
 * Događaj koji stakeholder servis objavljuje kada NE uspe da kreira profil.
 *
 * <p>Auth servis sluša ovaj događaj i pokreće kompenzaciju:
 * briše korisnika iz auth baze (rollback).
 *
 * <p>Ovo je ekvivalent UserRegisteredReply sa success=false iz orkestrisanog pristupa,
 * ali ovde je to DOGAĐAJ koji auth servis SAM detektuje i na koji reaguje.
 * Nema orchestratora koji koordiniše — svaki servis zna šta treba da uradi.
 *
 * <p>Routing key: "auth.user.rollback"
 * Exchange:      "saga.choreography.exchange"
 */
public class AuthUserRollbackEvent {

    private String sagaId;
    private String errorMessage;

    public AuthUserRollbackEvent() {}

    public AuthUserRollbackEvent(String sagaId, String errorMessage) {
        this.sagaId       = sagaId;
        this.errorMessage = errorMessage;
    }

    public String getSagaId()                  { return sagaId; }
    public void   setSagaId(String v)          { this.sagaId = v; }

    public String getErrorMessage()            { return errorMessage; }
    public void   setErrorMessage(String v)    { this.errorMessage = v; }
}
