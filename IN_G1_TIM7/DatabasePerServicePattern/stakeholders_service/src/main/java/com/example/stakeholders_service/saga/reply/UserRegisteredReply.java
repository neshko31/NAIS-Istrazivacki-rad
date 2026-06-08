package com.example.stakeholders_service.Saga.reply;

/**
 * Reply koji stakeholder servis šalje SagaOrchestratoru (auth servisu)
 * nakon obrade RegisterUserCommand.
 *
 * Polja:
 *   sagaId       – korelacioni ID sage
 *   success      – true ako je korisnik uspešno kreiran u stakeholder bazi
 *   errorMessage – opis greške (null ako je success=true)
 */
public class UserRegisteredReply {

    private String sagaId;
    private boolean success;
    private String errorMessage;

    public UserRegisteredReply() {}

    public UserRegisteredReply(String sagaId, boolean success, String errorMessage) {
        this.sagaId       = sagaId;
        this.success      = success;
        this.errorMessage = errorMessage;
    }

    public String  getSagaId()                    { return sagaId; }
    public void    setSagaId(String v)            { this.sagaId = v; }

    public boolean isSuccess()                    { return success; }
    public void    setSuccess(boolean v)          { this.success = v; }

    public String  getErrorMessage()              { return errorMessage; }
    public void    setErrorMessage(String v)      { this.errorMessage = v; }
}
