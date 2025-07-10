package tech.aiflowy.ai.voice.volcengine.asr.volcengine;

public class AsrResponse {
    private String reqid = "unknow";
    private int code = 0;
    private String message = "";
    private int sequence = 0;
    private Result[] result;
    private Addition addition;

    public String getReqid() {
        return reqid;
    }

    public void setReqid(String reqid) {
        this.reqid = reqid;
    }

    public int getCode() {
        return code;
    }

    public void setCode(int code) {
        this.code = code;
    }

    public String getMessage() {
        return message;
    }

    public void setMessage(String message) {
        this.message = message;
    }

    public int getSequence() {
        return sequence;
    }

    public void setSequence(int sequence) {
        this.sequence = sequence;
    }

    public Result[] getResult() {
        return result;
    }

    public void setResult(Result[] result) {
        this.result = result;
    }

    public Addition getAddition() {
        return addition;
    }

    public void setAddition(Addition addition) {
        this.addition = addition;
    }

    public static class Result {
        private String text;
        private int confidence;
        private String language;
        private Utterances[] utterances;
        private float global_confidence;

        public String getText() {
            return text;
        }

        public void setText(String text) {
            this.text = text;
        }

        public int getConfidence() {
            return confidence;
        }

        public void setConfidence(int confidence) {
            this.confidence = confidence;
        }

        public String getLanguage() {
            return language;
        }

        public void setLanguage(String language) {
            this.language = language;
        }

        public Utterances[] getUtterances() {
            return utterances;
        }

        public void setUtterances(Utterances[] utterances) {
            this.utterances = utterances;
        }

        public float getGlobal_confidence() {
            return global_confidence;
        }

        public void setGlobal_confidence(float global_confidence) {
            this.global_confidence = global_confidence;
        }
    }

    public static class Utterances {
        private String text;
        private int start_time;
        private int end_time;
        private boolean definite;
        private String language;
        private Words[] words;

        public String getText() {
            return text;
        }

        public void setText(String text) {
            this.text = text;
        }

        public int getStart_time() {
            return start_time;
        }

        public void setStart_time(int start_time) {
            this.start_time = start_time;
        }

        public int getEnd_time() {
            return end_time;
        }

        public void setEnd_time(int end_time) {
            this.end_time = end_time;
        }

        public boolean isDefinite() {
            return definite;
        }

        public void setDefinite(boolean definite) {
            this.definite = definite;
        }

        public String getLanguage() {
            return language;
        }

        public void setLanguage(String language) {
            this.language = language;
        }

        public Words[] getWords() {
            return words;
        }

        public void setWords(Words[] words) {
            this.words = words;
        }
    }

    public static class Words {
        private String text;
        private int start_time;
        private int end_time;
        private int blank_duration;

        public String getText() {
            return text;
        }

        public void setText(String text) {
            this.text = text;
        }

        public int getStart_time() {
            return start_time;
        }

        public void setStart_time(int start_time) {
            this.start_time = start_time;
        }

        public int getEnd_time() {
            return end_time;
        }

        public void setEnd_time(int end_time) {
            this.end_time = end_time;
        }

        public int getBlank_duration() {
            return blank_duration;
        }

        public void setBlank_duration(int blank_duration) {
            this.blank_duration = blank_duration;
        }
    }

    public static class Addition {
        private String duration;

        public String getDuration() {
            return duration;
        }

        public void setDuration(String duration) {
            this.duration = duration;
        }
    }
}