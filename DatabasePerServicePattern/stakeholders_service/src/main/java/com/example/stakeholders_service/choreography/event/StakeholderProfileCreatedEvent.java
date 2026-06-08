package com.example.stakeholders_service.choreography.event;

/**
 * Događaj koji stakeholder servis objavljuje nakon USPEŠNOG kreiranja profila.
 *
 * <p>Auth servis sluša ovaj događaj i označava sagu kao COMPLETED.
 * Ovo je ekvivalent UserRegisteredReply sa success=true iz orkestrisanog pristupa,
 * ali ovde je to DOGAĐAJ a ne reply na komandu.
 *
 * <p>Routing key: "stakeholder.profile.created"
 * Exchange:      "saga.choreography.exchange"
 */
public class StakeholderProfileCreatedEvent {

    private String  sagaId;
    private boolean success;

    public StakeholderProfileCreatedEvent() {}

    public StakeholderProfileCreatedEvent(String sagaId, boolean success) {
        this.sagaId  = sagaId;
        this.success = success;
    }

    public String  getSagaId()           { return sagaId; }
    public void    setSagaId(String v)   { this.sagaId = v; }

    public boolean isSuccess()           { return success; }
    public void    setSuccess(boolean v) { this.success = v; }
}
