package com.github.gigurra.kdump.internal.util

import scala.concurrent.duration.{Duration, FiniteDuration}
import scala.concurrent.{Await, ExecutionContext, Future}

object async {

  def run[T](expr: => T)(using ec: ExecutionContext): Future[T] =
    Future(expr)(ec)

  def await[T](f: Future[T])(using ec: ExecutionContext): T =
    Await.result(f, Duration.Inf)
}
