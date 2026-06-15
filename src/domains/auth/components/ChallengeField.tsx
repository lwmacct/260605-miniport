import { ReloadOutlined } from "@ant-design/icons";
import { Button, Input } from "antd";
import { useCallback, useEffect, useState } from "react";
import {
  createImageChallenge,
  type AuthChallengeConfig,
  type AuthChallengeResponse,
  type ImageChallenge,
} from "../api";
import styles from "./AuthScreen.module.css";

interface ChallengeFieldProps {
  config: AuthChallengeConfig;
  disabled?: boolean;
  resetKey: number;
  onChange(challenge?: AuthChallengeResponse): void;
  onError(message: string): void;
}

export function ChallengeField({
  config,
  disabled,
  resetKey,
  onChange,
  onError,
}: ChallengeFieldProps) {
  if (config.provider !== "image") {
    return (
      <div className={styles.remoteChallenge}>
        {config.sitekey ? "远程验证码已配置" : "远程验证码缺少站点公钥"}
      </div>
    );
  }

  return (
    <ImageChallengeField
      disabled={disabled}
      onChange={onChange}
      onError={onError}
      resetKey={resetKey}
    />
  );
}

interface ImageChallengeFieldProps {
  disabled?: boolean;
  resetKey: number;
  onChange(challenge?: AuthChallengeResponse): void;
  onError(message: string): void;
}

function ImageChallengeField({
  disabled,
  resetKey,
  onChange,
  onError,
}: ImageChallengeFieldProps) {
  const [challenge, setChallenge] = useState<ImageChallenge>();
  const [answer, setAnswer] = useState("");
  const [loading, setLoading] = useState(false);

  const refresh = useCallback(async () => {
    setLoading(true);
    setAnswer("");
    onChange(undefined);
    onError("");
    try {
      const next = await createImageChallenge();
      setChallenge(next);
    } catch (error) {
      setChallenge(undefined);
      onError(error instanceof Error ? error.message : "验证码生成失败");
    } finally {
      setLoading(false);
    }
  }, [onChange, onError]);

  useEffect(() => {
    void refresh();
  }, [refresh, resetKey]);

  useEffect(() => {
    if (!challenge || answer.trim() === "") {
      onChange(undefined);
      return;
    }

    onChange({
      provider: "image",
      challengeId: challenge.challengeId,
      answer,
    });
  }, [answer, challenge, onChange]);

  return (
    <div className={styles.captchaRow}>
      <Input
        autoComplete="off"
        className={styles.captchaInput}
        disabled={disabled || loading}
        onChange={(event) => setAnswer(event.target.value)}
        placeholder="输入验证码"
        value={answer}
        variant="borderless"
      />
      <Button
        className={styles.captchaButton}
        disabled={disabled || loading}
        htmlType="button"
        onClick={() => void refresh()}
      >
        {challenge ? <img alt="验证码" src={challenge.image} /> : <ReloadOutlined />}
      </Button>
    </div>
  );
}
