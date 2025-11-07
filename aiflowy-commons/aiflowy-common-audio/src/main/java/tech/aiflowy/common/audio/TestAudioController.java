package tech.aiflowy.common.audio;

import cn.dev33.satoken.annotation.SaIgnore;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;
import tech.aiflowy.common.audio.core.AudioServiceManager;
import tech.aiflowy.common.web.controller.BaseController;

import javax.annotation.Resource;

@SaIgnore
@RequestMapping("/tts")
@RestController
public class TestAudioController extends BaseController {

    @Resource
    private AudioServiceManager manager;

    @GetMapping("/test2")
    public String test2() throws Exception {
        return "2";
    }

    @GetMapping("/test1")
    public String test1() throws Exception {

        return "1";
    }

    @GetMapping("/test")
    public String test() throws Exception {

        return "hello world";
    }
}
